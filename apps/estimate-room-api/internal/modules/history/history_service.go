package history

import (
	"context"
	"errors"

	historydto "github.com/master-bogdan/estimate-room-api/internal/modules/history/dto"
	historyrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/history/repositories"
)

var errHistoryNotImplemented = errors.New("history not implemented")

type HistoryService interface {
	ListMySessions(ctx context.Context, userID string, query historydto.PaginationQuery) (historydto.PaginatedResponse[historydto.SessionListItem], error)
	ListTeamSessions(ctx context.Context, teamID, userID string, query historydto.PaginationQuery) (historydto.PaginatedResponse[historydto.SessionListItem], error)
	GetRoomSummary(ctx context.Context, roomID, userID string) (historydto.RoomSummaryResponse, error)
}

type historyService struct {
	repo historyrepositories.HistoryRepository
}

func NewHistoryService(repo historyrepositories.HistoryRepository) HistoryService {
	return &historyService{repo: repo}
}

func (s *historyService) ListMySessions(
	ctx context.Context,
	userID string,
	query historydto.PaginationQuery,
) (historydto.PaginatedResponse[historydto.SessionListItem], error) {
	return historydto.NewPaginatedResponse[historydto.SessionListItem](query, nil, 0), errHistoryNotImplemented
}

func (s *historyService) ListTeamSessions(
	ctx context.Context,
	teamID, userID string,
	query historydto.PaginationQuery,
) (historydto.PaginatedResponse[historydto.SessionListItem], error) {
	return historydto.NewPaginatedResponse[historydto.SessionListItem](query, nil, 0), errHistoryNotImplemented
}

func (s *historyService) GetRoomSummary(
	ctx context.Context,
	roomID, userID string,
) (historydto.RoomSummaryResponse, error) {
	return historydto.RoomSummaryResponse{}, errHistoryNotImplemented
}
