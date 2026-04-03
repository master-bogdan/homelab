package teams

import (
	"context"
	"errors"
	"strings"

	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type CreatedTeamInvite struct {
	Invitation *invitesmodels.InvitationModel
	Token      string
}

type TeamsInviteService interface {
	CreateInvites(ctx context.Context, teamID, actorUserID string, emails []string) ([]CreatedTeamInvite, error)
}

type teamsInviteService struct {
	db             *bun.DB
	teamRepo       teamsrepositories.TeamRepository
	memberRepo     teamsrepositories.TeamMemberRepository
	userService    users.UsersService
	invitesService invites.InvitesService
}

func NewTeamsInviteService(
	db *bun.DB,
	teamRepo teamsrepositories.TeamRepository,
	memberRepo teamsrepositories.TeamMemberRepository,
	userService users.UsersService,
	invitesService invites.InvitesService,
) TeamsInviteService {
	return &teamsInviteService{
		db:             db,
		teamRepo:       teamRepo,
		memberRepo:     memberRepo,
		userService:    userService,
		invitesService: invitesService,
	}
}

func (s *teamsInviteService) CreateInvites(
	ctx context.Context,
	teamID, actorUserID string,
	emails []string,
) ([]CreatedTeamInvite, error) {
	normalizedEmails := normalizeInviteEmails(emails)
	if len(normalizedEmails) == 0 {
		return nil, apperrors.ErrBadRequest
	}

	type inviteTarget struct {
		userID string
		email  string
	}

	targets := make([]inviteTarget, 0, len(normalizedEmails))
	for _, email := range normalizedEmails {
		user, err := s.userService.FindByEmail(email)
		if err != nil {
			if errors.Is(err, apperrors.ErrUserNotFound) {
				return nil, apperrors.ErrBadRequest
			}

			return nil, err
		}

		targets = append(targets, inviteTarget{
			userID: user.UserID,
			email:  email,
		})
	}

	createdInvites := make([]CreatedTeamInvite, 0, len(targets))
	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		teamRepo := teamsrepositories.NewTeamRepository(tx)
		memberRepo := teamsrepositories.NewTeamMemberRepository(tx)

		if _, err := ensureTeamOwner(teamRepo, memberRepo, teamID, actorUserID); err != nil {
			return err
		}

		for _, target := range targets {
			_, err := memberRepo.FindByTeamAndUser(teamID, target.userID)
			switch {
			case err == nil:
				return apperrors.ErrConflict
			case err != nil && !errors.Is(err, apperrors.ErrNotFound):
				return err
			}

			invitedUserID := target.userID
			invitedEmail := target.email

			invitation, token, err := s.invitesService.CreateInvitationWithDB(ctx, tx, invites.CreateInvitationInput{
				Kind:            invitesmodels.InvitationKindTeamMember,
				TeamID:          &teamID,
				InvitedUserID:   &invitedUserID,
				InvitedEmail:    &invitedEmail,
				CreatedByUserID: actorUserID,
			})
			if err != nil {
				return err
			}

			createdInvites = append(createdInvites, CreatedTeamInvite{
				Invitation: invitation,
				Token:      token,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return createdInvites, nil
}

func normalizeInviteEmails(emails []string) []string {
	seen := make(map[string]struct{}, len(emails))
	normalizedEmails := make([]string, 0, len(emails))

	for _, email := range emails {
		normalizedEmail := strings.ToLower(strings.TrimSpace(email))
		if normalizedEmail == "" {
			continue
		}
		if _, exists := seen[normalizedEmail]; exists {
			continue
		}

		seen[normalizedEmail] = struct{}{}
		normalizedEmails = append(normalizedEmails, normalizedEmail)
	}

	return normalizedEmails
}
