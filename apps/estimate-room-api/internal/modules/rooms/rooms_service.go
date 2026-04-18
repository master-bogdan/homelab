package rooms

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/master-bogdan/estimate-room-api/internal/modules/gamification"
	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	roomsutils "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/utils"
	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	usersrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/users/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/metrics"
	"github.com/uptrace/bun"
)

type RoomsService interface {
	CreateRoom(ctx context.Context, input CreateRoomInput) (*CreateRoomResult, error)
	GetRoom(roomID string) (*roomsmodels.RoomsModel, error)
	ValidateUserRoomAccess(roomID, userID string) error
	UpdateRoom(roomID, userID string, input UpdateRoomInput) (*roomsmodels.RoomsModel, error)
}

type CreateRoomInput struct {
	Name            string
	Deck            roomsmodels.RoomDeck
	AdminUserID     string
	InviteTeamID    *string
	InviteEmails    []string
	CreateShareLink bool
}

type CreatedRoomInvitation struct {
	Invitation *invitesmodels.InvitationModel
	Token      string
}

type CreateRoomSkippedRecipient struct {
	UserID *string
	Email  *string
	Reason string
}

type CreateRoomResult struct {
	Room              *roomsmodels.RoomsModel
	EmailInvitations  []CreatedRoomInvitation
	ShareLink         *CreatedRoomInvitation
	SkippedRecipients []CreateRoomSkippedRecipient
}

type roomsService struct {
	db              *bun.DB
	roomsRepo       roomsrepositories.RoomsRepository
	participantRepo roomsrepositories.RoomParticipantRepository
	teamRepo        teamsrepositories.TeamRepository
	memberRepo      teamsrepositories.TeamMemberRepository
	userRepo        usersrepositories.UserRepository
	invitesService  invites.InvitesService
	rewardService   gamification.RoomRewardService
	logger          *slog.Logger
}

func NewRoomsService(
	db *bun.DB,
	roomsRepo roomsrepositories.RoomsRepository,
	participantRepo roomsrepositories.RoomParticipantRepository,
	teamRepo teamsrepositories.TeamRepository,
	memberRepo teamsrepositories.TeamMemberRepository,
	userRepo usersrepositories.UserRepository,
	invitesService invites.InvitesService,
	rewardService gamification.RoomRewardService,
) RoomsService {
	return &roomsService{
		db:              db,
		roomsRepo:       roomsRepo,
		participantRepo: participantRepo,
		teamRepo:        teamRepo,
		memberRepo:      memberRepo,
		userRepo:        userRepo,
		invitesService:  invitesService,
		rewardService:   rewardService,
		logger:          logger.L().With(slog.String("service", "rooms")),
	}
}

func (s *roomsService) CreateRoom(ctx context.Context, input CreateRoomInput) (*CreateRoomResult, error) {
	model := roomsmodels.RoomsModel{
		Name:        strings.TrimSpace(input.Name),
		Deck:        input.Deck,
		AdminUserID: input.AdminUserID,
	}

	if model.Deck.IsZero() {
		model.Deck = roomsmodels.DefaultRoomDeck()
	}

	model.Deck.Name = strings.TrimSpace(model.Deck.Name)
	model.Deck.Kind = strings.TrimSpace(model.Deck.Kind)
	values := make([]string, 0, len(model.Deck.Values))
	for _, value := range model.Deck.Values {
		trimmedValue := strings.TrimSpace(value)
		if trimmedValue != "" {
			values = append(values, trimmedValue)
		}
	}
	model.Deck.Values = values

	if !model.Deck.IsValid() {
		return nil, fmt.Errorf("%w: invalid deck", apperrors.ErrBadRequest)
	}

	invitePlan, err := s.planRoomInvitations(input.AdminUserID, input.InviteTeamID, input.InviteEmails)
	if err != nil {
		return nil, err
	}
	if teamID := normalizeOptionalStringValue(input.InviteTeamID); teamID != "" {
		model.TeamID = &teamID
	}

	code, err := roomsutils.GenerateRoomCode()
	if err != nil {
		return nil, err
	}

	model.Code = code

	result := &CreateRoomResult{
		EmailInvitations:  make([]CreatedRoomInvitation, 0, len(invitePlan.Emails)),
		SkippedRecipients: invitePlan.SkippedRecipients,
	}
	err = s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		roomRepo := roomsrepositories.NewRoomsRepository(tx)
		participantRepo := roomsrepositories.NewRoomParticipantRepository(tx)

		room, err := roomRepo.Create(ctx, &model)
		if err != nil {
			return err
		}

		_, err = participantRepo.Create(&roomsmodels.RoomParticipantModel{
			RoomParticipantID: uuid.NewString(),
			RoomID:            room.RoomID,
			UserID:            &room.AdminUserID,
			Role:              roomsmodels.RoomParticipantRoleAdmin,
		})
		if err != nil {
			return err
		}

		result.Room = room

		for _, email := range invitePlan.Emails {
			emailCopy := email
			invitation, token, err := s.invitesService.CreateInvitationWithDB(ctx, tx, invites.CreateInvitationInput{
				Kind:            invitesmodels.InvitationKindRoomEmail,
				RoomID:          &room.RoomID,
				InvitedEmail:    &emailCopy,
				CreatedByUserID: input.AdminUserID,
			})
			if err != nil {
				return err
			}

			result.EmailInvitations = append(result.EmailInvitations, CreatedRoomInvitation{
				Invitation: invitation,
				Token:      token,
			})
		}

		if input.CreateShareLink {
			invitation, token, err := s.invitesService.CreateInvitationWithDB(ctx, tx, invites.CreateInvitationInput{
				Kind:            invitesmodels.InvitationKindRoomLink,
				RoomID:          &room.RoomID,
				CreatedByUserID: input.AdminUserID,
			})
			if err != nil {
				return err
			}

			result.ShareLink = &CreatedRoomInvitation{
				Invitation: invitation,
				Token:      token,
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	metrics.RecordRoomLifecycle("created")
	logger.FromContext(ctx, s.logger).Info(roomsServiceLog("Room created"), "room_id", result.Room.RoomID, "admin_user_id", input.AdminUserID, "team_id", result.Room.TeamID)

	return result, nil
}

func (s *roomsService) GetRoom(roomID string) (*roomsmodels.RoomsModel, error) {
	return s.roomsRepo.FindByID(roomID)
}

func (s *roomsService) ValidateUserRoomAccess(roomID, userID string) error {
	participant, err := s.participantRepo.FindActiveByUserID(roomID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrForbidden
		}
		return err
	}

	if !participant.Role.IsValid() {
		return apperrors.ErrForbidden
	}

	return nil
}

type UpdateRoomInput struct {
	Name   *string
	Status *string
}

func (s *roomsService) UpdateRoom(roomID, userID string, input UpdateRoomInput) (*roomsmodels.RoomsModel, error) {
	room, err := s.ensureRoomAdmin(roomID, userID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, apperrors.ErrBadRequest
		}
		input.Name = &name
	}
	if input.Status != nil {
		status := strings.TrimSpace(*input.Status)
		input.Status = &status
	}

	if isTerminalRoomStatus(room.Status) {
		return nil, fmt.Errorf("%w: terminal rooms cannot be updated", apperrors.ErrBadRequest)
	}
	if roomPatchIsNoop(room, input) {
		return room, nil
	}
	if input.Status != nil && *input.Status == "EXPIRED" {
		return nil, fmt.Errorf("%w: room expiry is system managed", apperrors.ErrBadRequest)
	}

	var updatedRoom *roomsmodels.RoomsModel
	err = s.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		roomRepo := roomsrepositories.NewRoomsRepository(tx)

		updatedRoom, err = roomRepo.Update(roomID, roomsrepositories.UpdateRoomFields{
			Name:   input.Name,
			Status: input.Status,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if room.Status != updatedRoom.Status && isTerminalRoomStatus(updatedRoom.Status) {
		metrics.RecordRoomLifecycle(updatedRoom.Status)
	}

	if s.rewardService != nil && isTerminalRoomStatus(updatedRoom.Status) {
		appliedRewards := s.applyTerminalRewardsBestEffort(updatedRoom)
		if len(appliedRewards) > 0 {
			if err := s.rewardService.NotifyAppliedRewards(context.Background(), appliedRewards); err != nil {
				s.logger.Error(roomsServiceLog("Failed to notify room rewards"), "room_id", roomID, "err", err)
			}
		}
	}

	s.logger.Info(roomsServiceLog("Room updated"), "room_id", updatedRoom.RoomID, "status", updatedRoom.Status, "admin_user_id", userID)

	return updatedRoom, nil
}

func (s *roomsService) applyTerminalRewardsBestEffort(room *roomsmodels.RoomsModel) []gamification.AppliedRoomReward {
	if s.rewardService == nil || room == nil || !isTerminalRoomStatus(room.Status) {
		return nil
	}

	appliedRewards := make([]gamification.AppliedRoomReward, 0)
	err := s.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		rewards, err := s.rewardService.ApplyRoomTerminalRewards(ctx, tx, room)
		if err != nil {
			return err
		}

		appliedRewards = rewards
		return nil
	})
	if err != nil {
		s.logger.Error(roomsServiceLog("Failed to apply room rewards"), "room_id", room.RoomID, "status", room.Status, "err", err)
		return nil
	}

	return appliedRewards
}

func roomsServiceLog(message string) string {
	return logger.Prefix("MODULE", "ROOMS", message)
}

func roomPatchIsNoop(room *roomsmodels.RoomsModel, input UpdateRoomInput) bool {
	if input.Name != nil && room.Name != *input.Name {
		return false
	}
	if input.Status != nil && room.Status != *input.Status {
		return false
	}
	return true
}

func isTerminalRoomStatus(status string) bool {
	return status == "FINISHED" || status == "EXPIRED"
}

type roomInvitationPlan struct {
	Emails            []string
	SkippedRecipients []CreateRoomSkippedRecipient
}

func (s *roomsService) planRoomInvitations(
	adminUserID string,
	inviteTeamID *string,
	inviteEmails []string,
) (*roomInvitationPlan, error) {
	plan := &roomInvitationPlan{
		Emails:            make([]string, 0),
		SkippedRecipients: make([]CreateRoomSkippedRecipient, 0),
	}

	creator, err := s.userRepo.FindByID(adminUserID)
	if err != nil {
		return nil, err
	}

	var creatorEmail string
	if creator.Email != nil {
		creatorEmail = strings.ToLower(strings.TrimSpace(*creator.Email))
	}

	seenEmails := make(map[string]struct{}, len(inviteEmails))
	seenSkipped := make(map[string]struct{})

	if teamID := normalizeOptionalStringValue(inviteTeamID); teamID != "" {
		team, err := s.ensureTeamInviteOwner(teamID, adminUserID)
		if err != nil {
			return nil, err
		}

		for _, member := range team.Members {
			if member == nil {
				continue
			}

			if member.UserID == adminUserID {
				continue
			}

			if member.User == nil || member.User.Email == nil || strings.TrimSpace(*member.User.Email) == "" {
				userID := member.UserID
				s.addSkippedRecipient(plan, seenSkipped, CreateRoomSkippedRecipient{
					UserID: &userID,
					Reason: "missing_email",
				})
				continue
			}

			normalizedEmail := strings.ToLower(strings.TrimSpace(*member.User.Email))
			if normalizedEmail == "" {
				userID := member.UserID
				s.addSkippedRecipient(plan, seenSkipped, CreateRoomSkippedRecipient{
					UserID: &userID,
					Reason: "missing_email",
				})
				continue
			}
			if creatorEmail != "" && normalizedEmail == creatorEmail {
				email := normalizedEmail
				s.addSkippedRecipient(plan, seenSkipped, CreateRoomSkippedRecipient{
					UserID: &member.UserID,
					Email:  &email,
					Reason: "self",
				})
				continue
			}
			if _, exists := seenEmails[normalizedEmail]; exists {
				continue
			}

			seenEmails[normalizedEmail] = struct{}{}
			plan.Emails = append(plan.Emails, normalizedEmail)
		}
	}

	for _, email := range normalizeInviteEmails(inviteEmails) {
		if creatorEmail != "" && email == creatorEmail {
			emailCopy := email
			s.addSkippedRecipient(plan, seenSkipped, CreateRoomSkippedRecipient{
				Email:  &emailCopy,
				Reason: "self",
			})
			continue
		}
		if _, exists := seenEmails[email]; exists {
			continue
		}

		seenEmails[email] = struct{}{}
		plan.Emails = append(plan.Emails, email)
	}

	return plan, nil
}

func (s *roomsService) ensureTeamInviteOwner(teamID, actorUserID string) (*teamsmodels.TeamModel, error) {
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return nil, err
	}

	member, err := s.memberRepo.FindByTeamAndUser(teamID, actorUserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrForbidden
		}

		return nil, err
	}

	if member.Role != teamsmodels.TeamMemberRoleOwner || team.OwnerUserID != actorUserID {
		return nil, apperrors.ErrForbidden
	}

	return team, nil
}

func (s *roomsService) addSkippedRecipient(
	plan *roomInvitationPlan,
	seen map[string]struct{},
	recipient CreateRoomSkippedRecipient,
) {
	key := recipient.Reason + "|" + valueOrEmpty(recipient.UserID) + "|" + valueOrEmpty(recipient.Email)
	if _, exists := seen[key]; exists {
		return
	}

	seen[key] = struct{}{}
	plan.SkippedRecipients = append(plan.SkippedRecipients, recipient)
}

func normalizeInviteEmails(emails []string) []string {
	seen := make(map[string]struct{}, len(emails))
	normalized := make([]string, 0, len(emails))

	for _, email := range emails {
		trimmed := strings.ToLower(strings.TrimSpace(email))
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}

		seen[trimmed] = struct{}{}
		normalized = append(normalized, trimmed)
	}

	return normalized
}

func normalizeOptionalStringValue(value *string) string {
	if value == nil {
		return ""
	}

	return strings.TrimSpace(*value)
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func (s *roomsService) ensureRoomAdmin(roomID, userID string) (*roomsmodels.RoomsModel, error) {
	room, err := s.roomsRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	participant, err := s.participantRepo.FindActiveByUserID(roomID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrForbidden
		}
		return nil, err
	}

	if participant.Role != roomsmodels.RoomParticipantRoleAdmin {
		return nil, apperrors.ErrForbidden
	}

	if room.AdminUserID != userID {
		return nil, apperrors.ErrForbidden
	}

	return room, nil
}
