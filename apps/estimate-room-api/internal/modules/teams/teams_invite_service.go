package teams

import (
	"context"
	"errors"
	"strings"

	"github.com/master-bogdan/estimate-room-api/internal/modules/invites"
	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	invitesrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/invites/repositories"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/users"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

type CreatedTeamInvite struct {
	Invitation *invitesmodels.InvitationModel
	Token      string
}

type TeamsInviteService interface {
	CreateInvites(ctx context.Context, teamID, actorUserID string, emails []string) ([]CreatedTeamInvite, error)
}

type teamsInviteService struct {
	teamRepo       teamsrepositories.TeamRepository
	memberRepo     teamsrepositories.TeamMemberRepository
	invitationRepo invitesrepositories.InvitationRepository
	userService    users.UsersService
	invitesService invites.InvitesService
}

func NewTeamsInviteService(
	teamRepo teamsrepositories.TeamRepository,
	memberRepo teamsrepositories.TeamMemberRepository,
	invitationRepo invitesrepositories.InvitationRepository,
	userService users.UsersService,
	invitesService invites.InvitesService,
) TeamsInviteService {
	return &teamsInviteService{
		teamRepo:       teamRepo,
		memberRepo:     memberRepo,
		invitationRepo: invitationRepo,
		userService:    userService,
		invitesService: invitesService,
	}
}

func (s *teamsInviteService) CreateInvites(
	ctx context.Context,
	teamID, actorUserID string,
	emails []string,
) ([]CreatedTeamInvite, error) {
	_, err := ensureTeamOwner(s.teamRepo, s.memberRepo, teamID, actorUserID)
	if err != nil {
		return nil, err
	}

	normalizedEmails := normalizeInviteEmails(emails)
	if len(normalizedEmails) == 0 {
		return nil, apperrors.ErrBadRequest
	}

	createdInvites := make([]CreatedTeamInvite, 0, len(normalizedEmails))
	for _, email := range normalizedEmails {
		user, err := s.userService.FindByEmail(email)
		if err != nil {
			if errors.Is(err, apperrors.ErrUserNotFound) {
				return nil, apperrors.ErrBadRequest
			}

			return nil, err
		}

		_, err = s.memberRepo.FindByTeamAndUser(teamID, user.UserID)
		switch {
		case err == nil:
			return nil, apperrors.ErrConflict
		case err != nil && !errors.Is(err, apperrors.ErrNotFound):
			return nil, err
		}

		_, err = s.invitationRepo.FindActiveTeamMemberInvitation(teamID, user.UserID)
		switch {
		case err == nil:
			return nil, apperrors.ErrConflict
		case err != nil && !errors.Is(err, apperrors.ErrNotFound):
			return nil, err
		}

		invitation, token, err := s.invitesService.CreateInvitation(ctx, invites.CreateInvitationInput{
			Kind:            invitesmodels.InvitationKindTeamMember,
			TeamID:          &teamID,
			InvitedUserID:   &user.UserID,
			InvitedEmail:    user.Email,
			CreatedByUserID: actorUserID,
		})
		if err != nil {
			return nil, err
		}

		createdInvites = append(createdInvites, CreatedTeamInvite{
			Invitation: invitation,
			Token:      token,
		})
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
