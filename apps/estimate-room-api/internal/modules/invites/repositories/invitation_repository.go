package invitesrepositories

import (
	"context"
	"database/sql"
	"errors"

	invitesmodels "github.com/master-bogdan/estimate-room-api/internal/modules/invites/models"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type InvitationRepository interface {
	Create(ctx context.Context, model *invitesmodels.InvitationModel) (*invitesmodels.InvitationModel, error)
	FindByID(invitationID string) (*invitesmodels.InvitationModel, error)
	FindByTokenID(tokenID string) (*invitesmodels.InvitationModel, error)
	FindActiveTeamMemberInvitation(teamID, invitedUserID string) (*invitesmodels.InvitationModel, error)
	Accept(invitationID string) (*invitesmodels.InvitationModel, error)
	Decline(invitationID string) (*invitesmodels.InvitationModel, error)
	Revoke(invitationID string) (*invitesmodels.InvitationModel, error)
}

type invitationRepository struct {
	db bun.IDB
}

func NewInvitationRepository(db bun.IDB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) Create(ctx context.Context, model *invitesmodels.InvitationModel) (*invitesmodels.InvitationModel, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	_, err := r.db.NewInsert().
		Model(model).
		Column(
			"invitation_id",
			"kind",
			"status",
			"team_id",
			"room_id",
			"invited_user_id",
			"invited_email",
			"created_by_user_id",
			"token_id",
		).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *invitationRepository) FindByID(invitationID string) (*invitesmodels.InvitationModel, error) {
	invitation := new(invitesmodels.InvitationModel)
	err := r.db.NewSelect().
		Model(invitation).
		Where("i.invitation_id = ?", invitationID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return invitation, nil
}

func (r *invitationRepository) FindByTokenID(tokenID string) (*invitesmodels.InvitationModel, error) {
	invitation := new(invitesmodels.InvitationModel)
	err := r.db.NewSelect().
		Model(invitation).
		Where("i.token_id = ?", tokenID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return invitation, nil
}

func (r *invitationRepository) FindActiveTeamMemberInvitation(
	teamID, invitedUserID string,
) (*invitesmodels.InvitationModel, error) {
	invitation := new(invitesmodels.InvitationModel)
	err := r.db.NewSelect().
		Model(invitation).
		Where("i.kind = ?", invitesmodels.InvitationKindTeamMember).
		Where("i.status = ?", invitesmodels.InvitationStatusActive).
		Where("i.team_id = ?", teamID).
		Where("i.invited_user_id = ?", invitedUserID).
		Limit(1).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return invitation, nil
}

func (r *invitationRepository) Accept(invitationID string) (*invitesmodels.InvitationModel, error) {
	return r.transition(invitationID, invitesmodels.InvitationStatusAccepted, "accepted_at")
}

func (r *invitationRepository) Decline(invitationID string) (*invitesmodels.InvitationModel, error) {
	return r.transition(invitationID, invitesmodels.InvitationStatusDeclined, "declined_at")
}

func (r *invitationRepository) Revoke(invitationID string) (*invitesmodels.InvitationModel, error) {
	return r.transition(invitationID, invitesmodels.InvitationStatusRevoked, "revoked_at")
}

func (r *invitationRepository) transition(
	invitationID string,
	status invitesmodels.InvitationStatus,
	timestampColumn string,
) (*invitesmodels.InvitationModel, error) {
	result, err := r.db.NewUpdate().
		Model((*invitesmodels.InvitationModel)(nil)).
		Set("status = ?", status).
		Set("updated_at = NOW()").
		Set(timestampColumn+" = NOW()").
		Where("invitation_id = ?", invitationID).
		Where("status = ?", invitesmodels.InvitationStatusActive).
		Exec(context.Background())
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, apperrors.ErrConflict
	}

	return r.FindByID(invitationID)
}
