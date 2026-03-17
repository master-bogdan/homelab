package invites

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
)

type CreateInvitationInput struct {
	Kind            invitesmodels.InvitationKind
	TeamID          *string
	RoomID          *string
	InvitedUserID   *string
	InvitedEmail    *string
	CreatedByUserID string
}

type InvitationTokenClaims struct {
	InvitationID string                       `json:"invitationId"`
	TokenID      string                       `json:"tokenId"`
	Kind         invitesmodels.InvitationKind `json:"kind"`
	TeamID       *string                      `json:"teamId,omitempty"`
	RoomID       *string                      `json:"roomId,omitempty"`
}

type roomGuestTokenClaims struct {
	RoomID        string                          `json:"roomId"`
	ParticipantID string                          `json:"participantId"`
	Role          roomsmodels.RoomParticipantRole `json:"role"`
	ExpiresAt     time.Time                       `json:"expiresAt"`
}

type AcceptInvitationResult struct {
	Invitation  *invitesmodels.InvitationModel
	Room        *roomsmodels.RoomsModel
	Participant *roomsmodels.RoomParticipantModel
	GuestToken  string
}

const (
	GuestAccessCookieName = "room_guest_token"
	guestTokenTTL         = 30 * 24 * time.Hour
)

type InvitesService interface {
	CreateInvitation(ctx context.Context, input CreateInvitationInput) (*invitesmodels.InvitationModel, string, error)
	CreateInvitationWithDB(ctx context.Context, db bun.IDB, input CreateInvitationInput) (*invitesmodels.InvitationModel, string, error)
	ParseInvitationToken(token string) (*InvitationTokenClaims, error)
	PreviewInvitation(ctx context.Context, token string) (*invitesmodels.InvitationModel, error)
	AcceptInvitation(ctx context.Context, token, actorUserID string, guestName *string) (*AcceptInvitationResult, error)
	DeclineInvitation(ctx context.Context, token, actorUserID string) (*invitesmodels.InvitationModel, error)
	RevokeInvitation(ctx context.Context, invitationID, actorUserID string) (*invitesmodels.InvitationModel, error)
	ValidateGuestRoomAccess(roomID, guestToken string) (*roomsmodels.RoomParticipantModel, error)
}

type invitesService struct {
	db              *bun.DB
	invitationRepo  invitesrepositories.InvitationRepository
	roomsRepo       roomsrepositories.RoomsRepository
	participantRepo roomsrepositories.RoomParticipantRepository
	teamRepo        teamsrepositories.TeamRepository
	memberRepo      teamsrepositories.TeamMemberRepository
	userRepo        usersrepositories.UserRepository
	tokenKey        []byte
}

func NewInvitesService(
	db *bun.DB,
	invitationRepo invitesrepositories.InvitationRepository,
	tokenKey string,
) InvitesService {
	return &invitesService{
		db:              db,
		invitationRepo:  invitationRepo,
		roomsRepo:       roomsrepositories.NewRoomsRepository(db),
		participantRepo: roomsrepositories.NewRoomParticipantRepository(db),
		teamRepo:        teamsrepositories.NewTeamRepository(db),
		memberRepo:      teamsrepositories.NewTeamMemberRepository(db),
		userRepo:        usersrepositories.NewUserRepository(db),
		tokenKey:        []byte(tokenKey),
	}
}

func (s *invitesService) CreateInvitation(
	ctx context.Context,
	input CreateInvitationInput,
) (*invitesmodels.InvitationModel, string, error) {
	return s.createInvitation(ctx, s.invitationRepo, input)
}

func (s *invitesService) CreateInvitationWithDB(
	ctx context.Context,
	db bun.IDB,
	input CreateInvitationInput,
) (*invitesmodels.InvitationModel, string, error) {
	if db == nil {
		return nil, "", apperrors.ErrInternal
	}

	return s.createInvitation(ctx, invitesrepositories.NewInvitationRepository(db), input)
}

func (s *invitesService) createInvitation(
	ctx context.Context,
	repo invitesrepositories.InvitationRepository,
	input CreateInvitationInput,
) (*invitesmodels.InvitationModel, string, error) {
	normalizedInput, err := normalizeCreateInvitationInput(input)
	if err != nil {
		return nil, "", err
	}

	invitation := &invitesmodels.InvitationModel{
		InvitationID:    uuid.NewString(),
		Kind:            normalizedInput.Kind,
		Status:          invitesmodels.InvitationStatusActive,
		TeamID:          normalizedInput.TeamID,
		RoomID:          normalizedInput.RoomID,
		InvitedUserID:   normalizedInput.InvitedUserID,
		InvitedEmail:    normalizedInput.InvitedEmail,
		CreatedByUserID: normalizedInput.CreatedByUserID,
		TokenID:         uuid.NewString(),
	}

	token, err := s.generateInvitationToken(invitation)
	if err != nil {
		return nil, "", err
	}

	createdInvitation, err := repo.Create(ctx, invitation)
	if err != nil {
		if isActiveTeamMemberInvitationConflict(err) {
			return nil, "", apperrors.ErrConflict
		}

		return nil, "", err
	}

	return createdInvitation, token, nil
}

func (s *invitesService) ParseInvitationToken(token string) (*InvitationTokenClaims, error) {
	trimmedToken := strings.TrimSpace(token)
	if trimmedToken == "" {
		return nil, apperrors.ErrNotFound
	}

	claims, err := utils.ParseToken[InvitationTokenClaims](s.tokenKey, trimmedToken)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	if claims.InvitationID == "" || claims.TokenID == "" || !claims.Kind.IsValid() {
		return nil, apperrors.ErrNotFound
	}

	return claims, nil
}

func (s *invitesService) PreviewInvitation(ctx context.Context, token string) (*invitesmodels.InvitationModel, error) {
	claims, err := s.ParseInvitationToken(token)
	if err != nil {
		return nil, err
	}

	invitation, err := s.invitationRepo.FindByTokenID(ctx, claims.TokenID)
	if err != nil {
		return nil, err
	}

	if !matchesInvitationClaims(invitation, claims) {
		return nil, apperrors.ErrNotFound
	}

	return invitation, nil
}

func (s *invitesService) AcceptInvitation(
	ctx context.Context,
	token, actorUserID string,
	guestName *string,
) (*AcceptInvitationResult, error) {
	invitation, err := s.PreviewInvitation(ctx, token)
	if err != nil {
		return nil, err
	}

	if invitation.Status != invitesmodels.InvitationStatusActive {
		return nil, apperrors.ErrConflict
	}

	switch invitation.Kind {
	case invitesmodels.InvitationKindTeamMember:
		acceptedInvitation, err := s.acceptTeamMemberInvitation(ctx, invitation, actorUserID)
		if err != nil {
			return nil, err
		}

		return &AcceptInvitationResult{Invitation: acceptedInvitation}, nil
	case invitesmodels.InvitationKindRoomEmail:
		return s.acceptRoomEmailInvitation(ctx, invitation, actorUserID)
	case invitesmodels.InvitationKindRoomLink:
		return s.acceptRoomLinkInvitation(ctx, invitation, actorUserID, guestName)
	default:
		return nil, apperrors.ErrBadRequest
	}
}

func (s *invitesService) DeclineInvitation(
	ctx context.Context,
	token, actorUserID string,
) (*invitesmodels.InvitationModel, error) {
	invitation, err := s.PreviewInvitation(ctx, token)
	if err != nil {
		return nil, err
	}

	if invitation.Status != invitesmodels.InvitationStatusActive {
		return nil, apperrors.ErrConflict
	}

	if invitation.Kind == invitesmodels.InvitationKindTeamMember {
		if err := s.ensureTeamMemberInviteActor(invitation, actorUserID); err != nil {
			return nil, err
		}
	}
	if invitation.Kind == invitesmodels.InvitationKindRoomEmail {
		if err := s.ensureRoomEmailInviteActor(invitation, actorUserID); err != nil {
			return nil, err
		}
	}
	if invitation.Kind == invitesmodels.InvitationKindRoomLink {
		return nil, apperrors.ErrBadRequest
	}

	return s.invitationRepo.Decline(ctx, invitation.InvitationID)
}

func (s *invitesService) RevokeInvitation(
	ctx context.Context,
	invitationID, actorUserID string,
) (*invitesmodels.InvitationModel, error) {
	invitation, err := s.invitationRepo.FindByID(ctx, strings.TrimSpace(invitationID))
	if err != nil {
		return nil, err
	}

	if invitation.Status != invitesmodels.InvitationStatusActive {
		return nil, apperrors.ErrConflict
	}

	if invitation.Kind == invitesmodels.InvitationKindTeamMember {
		if _, err := s.ensureTeamOwner(invitation, actorUserID); err != nil {
			return nil, err
		}
	}
	if invitation.Kind == invitesmodels.InvitationKindRoomEmail || invitation.Kind == invitesmodels.InvitationKindRoomLink {
		if _, err := s.ensureRoomAdmin(invitation, actorUserID); err != nil {
			return nil, err
		}
	}

	return s.invitationRepo.Revoke(ctx, invitation.InvitationID)
}

func (s *invitesService) ValidateGuestRoomAccess(roomID, guestToken string) (*roomsmodels.RoomParticipantModel, error) {
	claims, err := s.parseGuestToken(strings.TrimSpace(guestToken))
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	if claims.RoomID != roomID || claims.Role != roomsmodels.RoomParticipantRoleGuest {
		return nil, apperrors.ErrForbidden
	}

	if !claims.ExpiresAt.IsZero() && claims.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ErrUnauthorized
	}

	participant, err := s.participantRepo.FindActiveByID(roomID, claims.ParticipantID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrForbidden
		}

		return nil, err
	}

	if participant.Role != roomsmodels.RoomParticipantRoleGuest {
		return nil, apperrors.ErrForbidden
	}

	return participant, nil
}

func (s *invitesService) generateInvitationToken(invitation *invitesmodels.InvitationModel) (string, error) {
	return utils.GenerateToken(s.tokenKey, InvitationTokenClaims{
		InvitationID: invitation.InvitationID,
		TokenID:      invitation.TokenID,
		Kind:         invitation.Kind,
		TeamID:       invitation.TeamID,
		RoomID:       invitation.RoomID,
	})
}

func normalizeCreateInvitationInput(input CreateInvitationInput) (CreateInvitationInput, error) {
	input.CreatedByUserID = strings.TrimSpace(input.CreatedByUserID)
	if input.CreatedByUserID == "" {
		return CreateInvitationInput{}, apperrors.ErrBadRequest
	}

	input.TeamID = normalizeOptionalString(input.TeamID)
	input.RoomID = normalizeOptionalString(input.RoomID)
	input.InvitedUserID = normalizeOptionalString(input.InvitedUserID)
	input.InvitedEmail = normalizeOptionalEmail(input.InvitedEmail)

	if !input.Kind.IsValid() {
		return CreateInvitationInput{}, apperrors.ErrBadRequest
	}

	switch input.Kind {
	case invitesmodels.InvitationKindTeamMember:
		if input.TeamID == nil || input.RoomID != nil || input.InvitedUserID == nil || input.InvitedEmail == nil {
			return CreateInvitationInput{}, apperrors.ErrBadRequest
		}
	case invitesmodels.InvitationKindRoomEmail:
		if input.TeamID != nil || input.RoomID == nil || input.InvitedEmail == nil {
			return CreateInvitationInput{}, apperrors.ErrBadRequest
		}
	case invitesmodels.InvitationKindRoomLink:
		if input.TeamID != nil || input.RoomID == nil || input.InvitedUserID != nil || input.InvitedEmail != nil {
			return CreateInvitationInput{}, apperrors.ErrBadRequest
		}
	default:
		return CreateInvitationInput{}, apperrors.ErrBadRequest
	}

	return input, nil
}

func normalizeOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func normalizeOptionalEmail(value *string) *string {
	normalized := normalizeOptionalString(value)
	if normalized == nil {
		return nil
	}

	email := strings.ToLower(*normalized)
	return &email
}

func matchesInvitationClaims(
	invitation *invitesmodels.InvitationModel,
	claims *InvitationTokenClaims,
) bool {
	if invitation.InvitationID != claims.InvitationID {
		return false
	}
	if invitation.TokenID != claims.TokenID {
		return false
	}
	if invitation.Kind != claims.Kind {
		return false
	}
	if !sameOptionalString(invitation.TeamID, claims.TeamID) {
		return false
	}
	if !sameOptionalString(invitation.RoomID, claims.RoomID) {
		return false
	}

	return true
}

func sameOptionalString(left, right *string) bool {
	switch {
	case left == nil && right == nil:
		return true
	case left == nil || right == nil:
		return false
	default:
		return *left == *right
	}
}

func (s *invitesService) acceptRoomEmailInvitation(
	ctx context.Context,
	invitation *invitesmodels.InvitationModel,
	actorUserID string,
) (*AcceptInvitationResult, error) {
	if err := s.ensureRoomEmailInviteActor(invitation, actorUserID); err != nil {
		return nil, err
	}

	room, err := s.loadActiveRoomFromInvitation(invitation)
	if err != nil {
		return nil, err
	}

	var participant *roomsmodels.RoomParticipantModel
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		participantRepo := roomsrepositories.NewRoomParticipantRepository(tx)
		invitationRepo := invitesrepositories.NewInvitationRepository(tx)

		participant, err = participantRepo.FindActiveByUserID(room.RoomID, actorUserID)
		switch {
		case err == nil:
			_, err = invitationRepo.Accept(ctx, invitation.InvitationID)
			return err
		case err != nil && !errors.Is(err, apperrors.ErrNotFound):
			return err
		}

		participant, err = participantRepo.Create(&roomsmodels.RoomParticipantModel{
			RoomParticipantID: uuid.NewString(),
			RoomID:            room.RoomID,
			UserID:            &actorUserID,
			Role:              roomsmodels.RoomParticipantRoleMember,
		})
		if err != nil {
			return err
		}

		_, err = invitationRepo.Accept(ctx, invitation.InvitationID)
		return err
	})
	if err != nil {
		return nil, err
	}

	if err := s.roomsRepo.TouchActivity(room.RoomID); err != nil {
		return nil, err
	}

	fullRoom, err := s.roomsRepo.FindByID(room.RoomID)
	if err != nil {
		return nil, err
	}

	acceptedInvitation, err := s.invitationRepo.FindByID(ctx, invitation.InvitationID)
	if err != nil {
		return nil, err
	}

	return &AcceptInvitationResult{
		Invitation:  acceptedInvitation,
		Room:        fullRoom,
		Participant: participant,
	}, nil
}

func (s *invitesService) acceptRoomLinkInvitation(
	ctx context.Context,
	invitation *invitesmodels.InvitationModel,
	actorUserID string,
	guestName *string,
) (*AcceptInvitationResult, error) {
	room, err := s.loadActiveRoomFromInvitation(invitation)
	if err != nil {
		return nil, err
	}

	normalizedActorUserID := strings.TrimSpace(actorUserID)
	if normalizedActorUserID != "" {
		return s.joinRegisteredRoomParticipant(room, invitation, normalizedActorUserID)
	}

	if guestName == nil || strings.TrimSpace(*guestName) == "" {
		return nil, apperrors.ErrBadRequest
	}

	return s.joinGuestRoomParticipant(ctx, room, invitation, *guestName)
}

func (s *invitesService) acceptTeamMemberInvitation(
	ctx context.Context,
	invitation *invitesmodels.InvitationModel,
	actorUserID string,
) (*invitesmodels.InvitationModel, error) {
	if err := s.ensureTeamMemberInviteActor(invitation, actorUserID); err != nil {
		return nil, err
	}

	if invitation.TeamID == nil {
		return nil, apperrors.ErrBadRequest
	}

	_, err := s.teamRepo.FindByID(*invitation.TeamID)
	if err != nil {
		return nil, err
	}

	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		memberRepo := teamsrepositories.NewTeamMemberRepository(tx)
		invitationRepo := invitesrepositories.NewInvitationRepository(tx)

		_, err := memberRepo.FindByTeamAndUser(*invitation.TeamID, actorUserID)
		switch {
		case err == nil:
			_, err = invitationRepo.Accept(ctx, invitation.InvitationID)
			return err
		case err != nil && !errors.Is(err, apperrors.ErrNotFound):
			return err
		}

		_, err = memberRepo.Create(ctx, &teamsmodels.TeamMemberModel{
			TeamID: *invitation.TeamID,
			UserID: actorUserID,
			Role:   teamsmodels.TeamMemberRoleMember,
		})
		if err != nil {
			return err
		}

		_, err = invitationRepo.Accept(ctx, invitation.InvitationID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return s.invitationRepo.FindByID(ctx, invitation.InvitationID)
}

func (s *invitesService) joinRegisteredRoomParticipant(
	room *roomsmodels.RoomsModel,
	invitation *invitesmodels.InvitationModel,
	actorUserID string,
) (*AcceptInvitationResult, error) {
	participant, err := s.participantRepo.FindActiveByUserID(room.RoomID, actorUserID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	if participant == nil {
		participant, err = s.participantRepo.Create(&roomsmodels.RoomParticipantModel{
			RoomParticipantID: uuid.NewString(),
			RoomID:            room.RoomID,
			UserID:            &actorUserID,
			Role:              roomsmodels.RoomParticipantRoleMember,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := s.roomsRepo.TouchActivity(room.RoomID); err != nil {
		return nil, err
	}

	fullRoom, err := s.roomsRepo.FindByID(room.RoomID)
	if err != nil {
		return nil, err
	}

	return &AcceptInvitationResult{
		Invitation:  invitation,
		Room:        fullRoom,
		Participant: participant,
	}, nil
}

func (s *invitesService) joinGuestRoomParticipant(
	ctx context.Context,
	room *roomsmodels.RoomsModel,
	invitation *invitesmodels.InvitationModel,
	guestName string,
) (*AcceptInvitationResult, error) {
	trimmedGuestName := strings.TrimSpace(guestName)
	if trimmedGuestName == "" {
		return nil, apperrors.ErrBadRequest
	}

	participant, err := s.participantRepo.FindActiveByGuestName(room.RoomID, trimmedGuestName)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}
	if participant != nil {
		return nil, apperrors.ErrConflict
	}

	participant, err = s.participantRepo.Create(&roomsmodels.RoomParticipantModel{
		RoomParticipantID: uuid.NewString(),
		RoomID:            room.RoomID,
		GuestName:         &trimmedGuestName,
		Role:              roomsmodels.RoomParticipantRoleGuest,
	})
	if err != nil {
		return nil, err
	}

	guestToken, err := s.generateGuestToken(room.RoomID, participant.RoomParticipantID)
	if err != nil {
		return nil, err
	}

	if err := s.roomsRepo.TouchActivity(room.RoomID); err != nil {
		return nil, err
	}

	fullRoom, err := s.roomsRepo.FindByID(room.RoomID)
	if err != nil {
		return nil, err
	}

	return &AcceptInvitationResult{
		Invitation:  invitation,
		Room:        fullRoom,
		Participant: participant,
		GuestToken:  guestToken,
	}, nil
}

func (s *invitesService) ensureTeamMemberInviteActor(
	invitation *invitesmodels.InvitationModel,
	actorUserID string,
) error {
	normalizedActorUserID := strings.TrimSpace(actorUserID)
	if normalizedActorUserID == "" {
		return apperrors.ErrUnauthorized
	}

	if invitation.InvitedUserID == nil || *invitation.InvitedUserID != normalizedActorUserID {
		return apperrors.ErrForbidden
	}

	return nil
}

func (s *invitesService) ensureRoomEmailInviteActor(
	invitation *invitesmodels.InvitationModel,
	actorUserID string,
) error {
	normalizedActorUserID := strings.TrimSpace(actorUserID)
	if normalizedActorUserID == "" {
		return apperrors.ErrUnauthorized
	}

	user, err := s.userRepo.FindByID(normalizedActorUserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUnauthorized
		}

		return err
	}

	if invitation.InvitedEmail == nil || user.Email == nil || !strings.EqualFold(*invitation.InvitedEmail, *user.Email) {
		return apperrors.ErrForbidden
	}

	return nil
}

func (s *invitesService) ensureTeamOwner(
	invitation *invitesmodels.InvitationModel,
	actorUserID string,
) (*teamsmodels.TeamModel, error) {
	normalizedActorUserID := strings.TrimSpace(actorUserID)
	if normalizedActorUserID == "" {
		return nil, apperrors.ErrUnauthorized
	}

	if invitation.TeamID == nil {
		return nil, apperrors.ErrBadRequest
	}

	team, err := s.teamRepo.FindByID(*invitation.TeamID)
	if err != nil {
		return nil, err
	}

	member, err := s.memberRepo.FindByTeamAndUser(*invitation.TeamID, normalizedActorUserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrForbidden
		}

		return nil, err
	}

	if member.Role != teamsmodels.TeamMemberRoleOwner || team.OwnerUserID != normalizedActorUserID {
		return nil, apperrors.ErrForbidden
	}

	return team, nil
}

func (s *invitesService) ensureRoomAdmin(
	invitation *invitesmodels.InvitationModel,
	actorUserID string,
) (*roomsmodels.RoomsModel, error) {
	normalizedActorUserID := strings.TrimSpace(actorUserID)
	if normalizedActorUserID == "" {
		return nil, apperrors.ErrUnauthorized
	}

	if invitation.RoomID == nil {
		return nil, apperrors.ErrBadRequest
	}

	room, err := s.roomsRepo.FindByID(*invitation.RoomID)
	if err != nil {
		return nil, err
	}

	if room.AdminUserID != normalizedActorUserID {
		return nil, apperrors.ErrForbidden
	}

	return room, nil
}

func (s *invitesService) loadActiveRoomFromInvitation(
	invitation *invitesmodels.InvitationModel,
) (*roomsmodels.RoomsModel, error) {
	if invitation.RoomID == nil {
		return nil, apperrors.ErrBadRequest
	}

	room, err := s.roomsRepo.FindByID(*invitation.RoomID)
	if err != nil {
		return nil, err
	}

	if room.Status != "ACTIVE" {
		return nil, apperrors.ErrForbidden
	}

	return room, nil
}

func (s *invitesService) generateGuestToken(roomID, participantID string) (string, error) {
	return utils.GenerateToken(s.tokenKey, roomGuestTokenClaims{
		RoomID:        roomID,
		ParticipantID: participantID,
		Role:          roomsmodels.RoomParticipantRoleGuest,
		ExpiresAt:     time.Now().Add(guestTokenTTL),
	})
}

func (s *invitesService) parseGuestToken(token string) (*roomGuestTokenClaims, error) {
	return utils.ParseToken[roomGuestTokenClaims](s.tokenKey, token)
}

func isActiveTeamMemberInvitationConflict(err error) bool {
	var pgErr pgdriver.Error
	if !errors.As(err, &pgErr) {
		return false
	}

	if pgErr.Field('C') != pgerrcode.UniqueViolation {
		return false
	}

	return pgErr.Field('n') == "invitations_active_team_member_unique_idx"
}
