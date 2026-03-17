package invites

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
	"github.com/uptrace/bun"
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

type InvitesService interface {
	CreateInvitation(ctx context.Context, input CreateInvitationInput) (*invitesmodels.InvitationModel, string, error)
	ParseInvitationToken(token string) (*InvitationTokenClaims, error)
	PreviewInvitation(token string) (*invitesmodels.InvitationModel, error)
	AcceptInvitation(ctx context.Context, token, actorUserID string) (*invitesmodels.InvitationModel, error)
	DeclineInvitation(ctx context.Context, token, actorUserID string) (*invitesmodels.InvitationModel, error)
	RevokeInvitation(ctx context.Context, invitationID, actorUserID string) (*invitesmodels.InvitationModel, error)
}

type invitesService struct {
	db             *bun.DB
	invitationRepo invitesrepositories.InvitationRepository
	teamRepo       teamsrepositories.TeamRepository
	memberRepo     teamsrepositories.TeamMemberRepository
	tokenKey       []byte
}

func NewInvitesService(
	db *bun.DB,
	invitationRepo invitesrepositories.InvitationRepository,
	tokenKey string,
) InvitesService {
	return &invitesService{
		db:             db,
		invitationRepo: invitationRepo,
		teamRepo:       teamsrepositories.NewTeamRepository(db),
		memberRepo:     teamsrepositories.NewTeamMemberRepository(db),
		tokenKey:       []byte(tokenKey),
	}
}

func (s *invitesService) CreateInvitation(
	ctx context.Context,
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

	createdInvitation, err := s.invitationRepo.Create(ctx, invitation)
	if err != nil {
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

func (s *invitesService) PreviewInvitation(token string) (*invitesmodels.InvitationModel, error) {
	claims, err := s.ParseInvitationToken(token)
	if err != nil {
		return nil, err
	}

	invitation, err := s.invitationRepo.FindByTokenID(claims.TokenID)
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
) (*invitesmodels.InvitationModel, error) {
	invitation, err := s.PreviewInvitation(token)
	if err != nil {
		return nil, err
	}

	if invitation.Status != invitesmodels.InvitationStatusActive {
		return nil, apperrors.ErrConflict
	}

	if invitation.Kind == invitesmodels.InvitationKindTeamMember {
		return s.acceptTeamMemberInvitation(ctx, invitation, actorUserID)
	}

	return s.invitationRepo.Accept(invitation.InvitationID)
}

func (s *invitesService) DeclineInvitation(
	ctx context.Context,
	token, actorUserID string,
) (*invitesmodels.InvitationModel, error) {
	invitation, err := s.PreviewInvitation(token)
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

	return s.invitationRepo.Decline(invitation.InvitationID)
}

func (s *invitesService) RevokeInvitation(
	ctx context.Context,
	invitationID, actorUserID string,
) (*invitesmodels.InvitationModel, error) {
	invitation, err := s.invitationRepo.FindByID(strings.TrimSpace(invitationID))
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

	return s.invitationRepo.Revoke(invitation.InvitationID)
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
			_, err = invitationRepo.Accept(invitation.InvitationID)
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

		_, err = invitationRepo.Accept(invitation.InvitationID)
		return err
	})
	if err != nil {
		return nil, err
	}

	return s.invitationRepo.FindByID(invitation.InvitationID)
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
