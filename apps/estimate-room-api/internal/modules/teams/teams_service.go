package teams

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/uptrace/bun"
)

type TeamsService interface {
	CreateTeam(ctx context.Context, name, ownerUserID string) (*teamsmodels.TeamModel, error)
	ListTeams(userID string) ([]*teamsmodels.TeamModel, error)
	GetTeam(teamID, userID string) (*teamsmodels.TeamModel, error)
	RemoveMember(teamID, actorUserID, targetUserID string) error
}

type teamsService struct {
	db         *bun.DB
	teamRepo   teamsrepositories.TeamRepository
	memberRepo teamsrepositories.TeamMemberRepository
	logger     *slog.Logger
}

func NewTeamsService(
	db *bun.DB,
	teamRepo teamsrepositories.TeamRepository,
	memberRepo teamsrepositories.TeamMemberRepository,
) TeamsService {
	return &teamsService{
		db:         db,
		teamRepo:   teamRepo,
		memberRepo: memberRepo,
		logger:     logger.L().With(slog.String("service", "teams")),
	}
}

func (s *teamsService) CreateTeam(ctx context.Context, name, ownerUserID string) (*teamsmodels.TeamModel, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, apperrors.ErrBadRequest
	}

	teamID := uuid.NewString()

	err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		teamRepo := teamsrepositories.NewTeamRepository(tx)
		memberRepo := teamsrepositories.NewTeamMemberRepository(tx)

		_, err := teamRepo.Create(ctx, &teamsmodels.TeamModel{
			TeamID:      teamID,
			Name:        trimmedName,
			OwnerUserID: ownerUserID,
		})
		if err != nil {
			return err
		}

		_, err = memberRepo.Create(ctx, &teamsmodels.TeamMemberModel{
			TeamID: teamID,
			UserID: ownerUserID,
			Role:   teamsmodels.TeamMemberRoleOwner,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.teamRepo.FindByID(teamID)
}

func (s *teamsService) ListTeams(userID string) ([]*teamsmodels.TeamModel, error) {
	return s.teamRepo.ListByUserID(userID)
}

func (s *teamsService) GetTeam(teamID, userID string) (*teamsmodels.TeamModel, error) {
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return nil, err
	}

	_, err = s.memberRepo.FindByTeamAndUser(teamID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrForbidden
		}

		return nil, err
	}

	return team, nil
}

func (s *teamsService) RemoveMember(teamID, actorUserID, targetUserID string) error {
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return err
	}

	actorMember, err := s.memberRepo.FindByTeamAndUser(teamID, actorUserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrForbidden
		}

		return err
	}

	if actorMember.Role != teamsmodels.TeamMemberRoleOwner || team.OwnerUserID != actorUserID {
		return apperrors.ErrForbidden
	}

	targetMember, err := s.memberRepo.FindByTeamAndUser(teamID, targetUserID)
	if err != nil {
		return err
	}

	if targetMember.Role == teamsmodels.TeamMemberRoleOwner || targetUserID == team.OwnerUserID {
		return apperrors.ErrBadRequest
	}

	return s.memberRepo.Delete(teamID, targetUserID)
}
