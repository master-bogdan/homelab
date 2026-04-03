package history

import (
	"context"
	"errors"
	"strings"

	historydto "github.com/master-bogdan/estimate-room-api/internal/modules/history/dto"
	historyrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/history/repositories"
	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

var errHistoryNotImplemented = errors.New("history not implemented")

type HistoryService interface {
	ListMySessions(ctx context.Context, userID string, query historydto.MySessionsQuery) (historydto.PaginatedResponse[historydto.SessionListItem], error)
	ListTeamSessions(ctx context.Context, teamID, userID string, query historydto.TeamSessionsQuery) (historydto.PaginatedResponse[historydto.SessionListItem], error)
	GetRoomSummary(ctx context.Context, roomID, userID string) (historydto.RoomSummaryResponse, error)
}

type historyService struct {
	repo       historyrepositories.HistoryRepository
	teamRepo   teamsrepositories.TeamRepository
	memberRepo teamsrepositories.TeamMemberRepository
}

func NewHistoryService(
	repo historyrepositories.HistoryRepository,
	teamRepo teamsrepositories.TeamRepository,
	memberRepo teamsrepositories.TeamMemberRepository,
) HistoryService {
	return &historyService{
		repo:       repo,
		teamRepo:   teamRepo,
		memberRepo: memberRepo,
	}
}

func (s *historyService) ListMySessions(
	ctx context.Context,
	userID string,
	query historydto.MySessionsQuery,
) (historydto.PaginatedResponse[historydto.SessionListItem], error) {
	if strings.TrimSpace(userID) == "" {
		return historydto.NewPaginatedResponse[historydto.SessionListItem](query.PaginationQuery, nil, 0), apperrors.ErrBadRequest
	}

	items, total, err := s.repo.ListMySessions(ctx, userID, query)
	if err != nil {
		return historydto.NewPaginatedResponse[historydto.SessionListItem](query.PaginationQuery, nil, 0), err
	}

	return historydto.NewPaginatedResponse(query.PaginationQuery, items, total), nil
}

func (s *historyService) ListTeamSessions(
	ctx context.Context,
	teamID, userID string,
	query historydto.TeamSessionsQuery,
) (historydto.PaginatedResponse[historydto.SessionListItem], error) {
	if strings.TrimSpace(teamID) == "" || strings.TrimSpace(userID) == "" {
		return historydto.NewPaginatedResponse[historydto.SessionListItem](query.PaginationQuery, nil, 0), apperrors.ErrBadRequest
	}

	if err := s.ensureTeamOwner(teamID, userID); err != nil {
		return historydto.NewPaginatedResponse[historydto.SessionListItem](query.PaginationQuery, nil, 0), err
	}

	items, total, err := s.repo.ListTeamSessions(ctx, teamID, userID, query)
	if err != nil {
		return historydto.NewPaginatedResponse[historydto.SessionListItem](query.PaginationQuery, nil, 0), err
	}

	return historydto.NewPaginatedResponse(query.PaginationQuery, items, total), nil
}

func (s *historyService) GetRoomSummary(
	ctx context.Context,
	roomID, userID string,
) (historydto.RoomSummaryResponse, error) {
	if strings.TrimSpace(roomID) == "" || strings.TrimSpace(userID) == "" {
		return historydto.RoomSummaryResponse{}, apperrors.ErrBadRequest
	}

	summary, err := s.repo.GetRoomSummary(ctx, roomID)
	if err != nil {
		return historydto.RoomSummaryResponse{}, err
	}

	if summary.Overview.AdminUser.UserID == userID {
		return summary, nil
	}

	if summary.Overview.TeamID == nil {
		return historydto.RoomSummaryResponse{}, apperrors.ErrForbidden
	}

	if err := s.ensureTeamOwner(*summary.Overview.TeamID, userID); err != nil {
		return historydto.RoomSummaryResponse{}, err
	}

	return summary, nil
}

func (s *historyService) ensureTeamOwner(teamID, userID string) error {
	team, err := s.teamRepo.FindByID(teamID)
	if err != nil {
		return err
	}

	member, err := s.memberRepo.FindByTeamAndUser(teamID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrForbidden
		}

		return err
	}

	if member.Role != teamsmodels.TeamMemberRoleOwner || team.OwnerUserID != userID {
		return apperrors.ErrForbidden
	}

	return nil
}
