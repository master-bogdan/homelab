package teams

import (
	"errors"

	teamsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/teams/models"
	teamsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/teams/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
)

func ensureTeamOwner(
	teamRepo teamsrepositories.TeamRepository,
	memberRepo teamsrepositories.TeamMemberRepository,
	teamID, actorUserID string,
) (*teamsmodels.TeamModel, error) {
	team, err := teamRepo.FindByID(teamID)
	if err != nil {
		return nil, err
	}

	actorMember, err := memberRepo.FindByTeamAndUser(teamID, actorUserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrForbidden
		}

		return nil, err
	}

	if actorMember.Role != teamsmodels.TeamMemberRoleOwner || team.OwnerUserID != actorUserID {
		return nil, apperrors.ErrForbidden
	}

	return team, nil
}
