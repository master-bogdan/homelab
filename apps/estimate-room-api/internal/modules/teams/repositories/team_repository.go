package teamsrepositories

import (
	"context"
	"database/sql"
	"errors"

	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type TeamRepository interface {
	Create(ctx context.Context, model *teamsmodels.TeamModel) (*teamsmodels.TeamModel, error)
	FindByID(teamID string) (*teamsmodels.TeamModel, error)
	ListByUserID(userID string) ([]*teamsmodels.TeamModel, error)
}

type teamRepository struct {
	db bun.IDB
}

func NewTeamRepository(db bun.IDB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, model *teamsmodels.TeamModel) (*teamsmodels.TeamModel, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column("team_id", "name", "owner_user_id").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *teamRepository) FindByID(teamID string) (*teamsmodels.TeamModel, error) {
	team := new(teamsmodels.TeamModel)
	err := r.db.NewSelect().
		Model(team).
		Where("t.team_id = ?", teamID).
		Relation("Members", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.
				OrderExpr("tm.joined_at ASC").
				Relation("User")
		}).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return team, nil
}

func (r *teamRepository) ListByUserID(userID string) ([]*teamsmodels.TeamModel, error) {
	teams := make([]*teamsmodels.TeamModel, 0)
	err := r.db.NewSelect().
		Model(&teams).
		Join("JOIN team_members AS tm ON tm.team_id = t.team_id").
		Where("tm.user_id = ?", userID).
		OrderExpr("t.created_at DESC").
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return teams, nil
}
