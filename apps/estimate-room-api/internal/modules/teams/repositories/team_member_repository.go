package teamsrepositories

import (
	"context"
	"database/sql"
	"errors"

	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type TeamMemberRepository interface {
	Create(ctx context.Context, model *teamsmodels.TeamMemberModel) (*teamsmodels.TeamMemberModel, error)
	FindByTeamAndUser(teamID, userID string) (*teamsmodels.TeamMemberModel, error)
	Delete(teamID, userID string) error
}

type teamMemberRepository struct {
	db bun.IDB
}

func NewTeamMemberRepository(db bun.IDB) TeamMemberRepository {
	return &teamMemberRepository{db: db}
}

func (r *teamMemberRepository) Create(ctx context.Context, model *teamsmodels.TeamMemberModel) (*teamsmodels.TeamMemberModel, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column("team_id", "user_id", "role").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *teamMemberRepository) FindByTeamAndUser(teamID, userID string) (*teamsmodels.TeamMemberModel, error) {
	member := new(teamsmodels.TeamMemberModel)
	err := r.db.NewSelect().
		Model(member).
		Where("tm.team_id = ?", teamID).
		Where("tm.user_id = ?", userID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return member, nil
}

func (r *teamMemberRepository) Delete(teamID, userID string) error {
	result, err := r.db.NewDelete().
		Model((*teamsmodels.TeamMemberModel)(nil)).
		Where("team_id = ?", teamID).
		Where("user_id = ?", userID).
		Exec(context.Background())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
