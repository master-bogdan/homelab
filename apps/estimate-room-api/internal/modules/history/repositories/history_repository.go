package historyrepositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	historydto "github.com/master-bogdan/estimate-room-api/internal/modules/history/dto"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/uptrace/bun"
)

type HistoryRepository interface {
	ListMySessions(ctx context.Context, userID string, query historydto.MySessionsQuery) ([]historydto.SessionListItem, int, error)
	ListTeamSessions(ctx context.Context, teamID, userID string, query historydto.TeamSessionsQuery) ([]historydto.SessionListItem, int, error)
	GetRoomSummary(ctx context.Context, roomID string) (historydto.RoomSummaryResponse, error)
}

type historyRepository struct {
	db *bun.DB
}

func NewHistoryRepository(db *bun.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) ListMySessions(
	ctx context.Context,
	userID string,
	query historydto.MySessionsQuery,
) ([]historydto.SessionListItem, int, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	userID = strings.TrimSpace(userID)
	whereSQL, whereArgs := buildListMySessionsWhere(userID, query)
	return r.listSessions(ctx, buildListSessionsInput{
		WhereSQL:   whereSQL,
		WhereArgs:  whereArgs,
		Query:      query.PaginationQuery,
		RoleUserID: userID,
	})
}

func (r *historyRepository) ListTeamSessions(
	ctx context.Context,
	teamID, userID string,
	query historydto.TeamSessionsQuery,
) ([]historydto.SessionListItem, int, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	teamID = strings.TrimSpace(teamID)
	userID = strings.TrimSpace(userID)
	whereSQL, whereArgs := buildListTeamSessionsWhere(teamID, query)
	return r.listSessions(ctx, buildListSessionsInput{
		WhereSQL:   whereSQL,
		WhereArgs:  whereArgs,
		Query:      query.PaginationQuery,
		RoleUserID: userID,
	})
}

func (r *historyRepository) GetRoomSummary(
	ctx context.Context,
	roomID string,
) (historydto.RoomSummaryResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return historydto.RoomSummaryResponse{}, apperrors.ErrBadRequest
	}

	overview, err := r.getRoomSummaryOverview(ctx, roomID)
	if err != nil {
		return historydto.RoomSummaryResponse{}, err
	}

	participants, err := r.getRoomSummaryParticipants(ctx, roomID)
	if err != nil {
		return historydto.RoomSummaryResponse{}, err
	}

	tasks, err := r.getRoomSummaryTasks(ctx, roomID)
	if err != nil {
		return historydto.RoomSummaryResponse{}, err
	}

	rounds, err := r.getRoomSummaryRounds(ctx, roomID)
	if err != nil {
		return historydto.RoomSummaryResponse{}, err
	}

	votes, err := r.getRoomSummaryVotes(ctx, roomID)
	if err != nil {
		return historydto.RoomSummaryResponse{}, err
	}

	roundsByTask := make(map[string][]historydto.RoomSummaryTaskRound, len(tasks))
	for _, round := range rounds {
		round.Votes = make([]historydto.RoomSummaryVote, 0)
		roundsByTask[round.TaskID] = append(roundsByTask[round.TaskID], round)
	}

	votesByTaskRound := make(map[string][]historydto.RoomSummaryVote, len(votes))
	for _, vote := range votes {
		key := roomTaskRoundKey(vote.TaskID, vote.RoundNumber)
		votesByTaskRound[key] = append(votesByTaskRound[key], vote)
	}

	for taskIdx := range tasks {
		taskRounds := roundsByTask[tasks[taskIdx].TaskID]
		for roundIdx := range taskRounds {
			key := roomTaskRoundKey(taskRounds[roundIdx].TaskID, taskRounds[roundIdx].RoundNumber)
			taskRounds[roundIdx].Votes = votesByTaskRound[key]
		}
		tasks[taskIdx].Rounds = taskRounds
	}

	return historydto.RoomSummaryResponse{
		Overview:     overview,
		Participants: participants,
		Tasks:        tasks,
	}, nil
}

type buildListSessionsInput struct {
	WhereSQL   string
	WhereArgs  []any
	Query      historydto.PaginationQuery
	RoleUserID string
}

func (r *historyRepository) listSessions(
	ctx context.Context,
	input buildListSessionsInput,
) ([]historydto.SessionListItem, int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM rooms AS r
		WHERE ` + input.WhereSQL

	if err := r.db.NewRaw(countQuery, input.WhereArgs...).Scan(ctx, &total); err != nil {
		return nil, 0, err
	}

	items := make([]historydto.SessionListItem, 0)
	listArgs := make([]any, 0, 2+len(input.WhereArgs)+2)
	listArgs = append(listArgs, input.RoleUserID, input.RoleUserID)
	listArgs = append(listArgs, input.WhereArgs...)
	listArgs = append(listArgs, input.Query.PageSize, input.Query.Offset())

	listQuery := `
		SELECT
			r.room_id,
			r.team_id,
			r.name,
			r.status,
			CASE
				WHEN r.admin_user_id = ? THEN 'ADMIN'
				WHEN EXISTS (
					SELECT 1
					FROM room_participants AS rp_viewer
					WHERE rp_viewer.room_id = r.room_id
					  AND rp_viewer.user_id = ?
				) THEN 'PARTICIPANT'
				ELSE 'VIEWER'
			END AS role,
			r.created_at,
			r.finished_at,
			r.last_activity_at,
			GREATEST(
				EXTRACT(EPOCH FROM (COALESCE(r.finished_at, r.last_activity_at) - r.created_at)),
				0
			)::bigint AS approx_duration_seconds,
			COALESCE((
				SELECT COUNT(*)
				FROM room_participants AS rp
				WHERE rp.room_id = r.room_id
			), 0)::int AS participants_count,
			COALESCE((
				SELECT COUNT(*)
				FROM tasks AS t
				WHERE t.room_id = r.room_id
			), 0)::int AS tasks_count,
			COALESCE((
				SELECT COUNT(*)
				FROM tasks AS t
				WHERE t.room_id = r.room_id
				  AND t.status = 'ESTIMATED'
			), 0)::int AS estimated_tasks_count
		FROM rooms AS r
		WHERE ` + input.WhereSQL + `
		ORDER BY COALESCE(r.finished_at, r.last_activity_at) DESC, r.created_at DESC, r.room_id DESC
		LIMIT ? OFFSET ?
	`

	if err := r.db.NewRaw(listQuery, listArgs...).Scan(ctx, &items); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *historyRepository) getRoomSummaryOverview(
	ctx context.Context,
	roomID string,
) (historydto.RoomSummaryOverview, error) {
	type roomSummaryOverviewRow struct {
		RoomID                string     `bun:"room_id"`
		TeamID                *string    `bun:"team_id"`
		Name                  string     `bun:"name"`
		Status                string     `bun:"status"`
		CreatedAt             time.Time  `bun:"created_at"`
		FinishedAt            *time.Time `bun:"finished_at"`
		LastActivityAt        time.Time  `bun:"last_activity_at"`
		ApproxDurationSeconds int64      `bun:"approx_duration_seconds"`
		ParticipantsCount     int        `bun:"participants_count"`
		EstimatedTasksCount   int        `bun:"estimated_tasks_count"`
		TasksCount            int        `bun:"tasks_count"`
		RoundCount            int        `bun:"round_count"`
		AdminUserID           string     `bun:"admin_user_id"`
		AdminEmail            *string    `bun:"admin_email"`
		AdminDisplayName      string     `bun:"admin_display_name"`
		AdminAvatarURL        *string    `bun:"admin_avatar_url"`
	}

	row := roomSummaryOverviewRow{}
	query := `
		SELECT
			r.room_id,
			r.team_id,
			r.name,
			r.status,
			r.created_at,
			r.finished_at,
			r.last_activity_at,
			GREATEST(
				EXTRACT(EPOCH FROM (COALESCE(r.finished_at, r.last_activity_at) - r.created_at)),
				0
			)::bigint AS approx_duration_seconds,
			COALESCE((
				SELECT COUNT(*)
				FROM room_participants AS rp
				WHERE rp.room_id = r.room_id
			), 0)::int AS participants_count,
			COALESCE((
				SELECT COUNT(*)
				FROM tasks AS t
				WHERE t.room_id = r.room_id
				  AND t.status = 'ESTIMATED'
			), 0)::int AS estimated_tasks_count,
			COALESCE((
				SELECT COUNT(*)
				FROM tasks AS t
				WHERE t.room_id = r.room_id
			), 0)::int AS tasks_count,
			COALESCE((
				SELECT COUNT(*)
				FROM task_rounds AS tr
				JOIN tasks AS t ON t.task_id = tr.task_id
				WHERE t.room_id = r.room_id
			), 0)::int AS round_count,
			r.admin_user_id,
			u.email AS admin_email,
			u.display_name AS admin_display_name,
			u.avatar_url AS admin_avatar_url
		FROM rooms AS r
		JOIN users AS u ON u.user_id = r.admin_user_id
		WHERE r.room_id = ?
		LIMIT 1
	`

	if err := r.db.NewRaw(query, roomID).Scan(ctx, &row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return historydto.RoomSummaryOverview{}, apperrors.ErrNotFound
		}
		return historydto.RoomSummaryOverview{}, err
	}

	return historydto.RoomSummaryOverview{
		RoomID:                row.RoomID,
		TeamID:                row.TeamID,
		Name:                  row.Name,
		Status:                row.Status,
		CreatedAt:             row.CreatedAt,
		FinishedAt:            row.FinishedAt,
		LastActivityAt:        row.LastActivityAt,
		ApproxDurationSeconds: row.ApproxDurationSeconds,
		ParticipantsCount:     row.ParticipantsCount,
		EstimatedTasksCount:   row.EstimatedTasksCount,
		TasksCount:            row.TasksCount,
		RoundCount:            row.RoundCount,
		AdminUser: historydto.RoomSummaryUserRef{
			UserID:      row.AdminUserID,
			Email:       row.AdminEmail,
			DisplayName: row.AdminDisplayName,
			AvatarURL:   row.AdminAvatarURL,
		},
	}, nil
}

func (r *historyRepository) getRoomSummaryParticipants(
	ctx context.Context,
	roomID string,
) ([]historydto.RoomSummaryParticipant, error) {
	participants := make([]historydto.RoomSummaryParticipant, 0)
	query := `
		SELECT
			rp.room_participants_id AS participant_id,
			rp.user_id,
			rp.guest_name,
			u.email,
			NULLIF(u.display_name, '') AS display_name,
			u.avatar_url,
			rp.role,
			rp.joined_at,
			rp.left_at,
			COALESCE(COUNT(v.votes_id), 0)::int AS votes_cast_count,
			COALESCE(COUNT(DISTINCT CASE WHEN t.status = 'ESTIMATED' THEN t.task_id END), 0)::int AS estimated_tasks_voted_count
		FROM room_participants AS rp
		LEFT JOIN users AS u ON u.user_id = rp.user_id
		LEFT JOIN votes AS v ON v.participant_id = rp.room_participants_id
		LEFT JOIN tasks AS t ON t.task_id = v.task_id
		WHERE rp.room_id = ?
		GROUP BY
			rp.room_participants_id,
			rp.user_id,
			rp.guest_name,
			u.email,
			u.display_name,
			u.avatar_url,
			rp.role,
			rp.joined_at,
			rp.left_at
		ORDER BY rp.joined_at ASC, rp.room_participants_id ASC
	`

	if err := r.db.NewRaw(query, roomID).Scan(ctx, &participants); err != nil {
		return nil, err
	}

	return participants, nil
}

func (r *historyRepository) getRoomSummaryTasks(
	ctx context.Context,
	roomID string,
) ([]historydto.RoomSummaryTask, error) {
	tasks := make([]historydto.RoomSummaryTask, 0)
	query := `
		SELECT
			t.task_id,
			t.title,
			t.description,
			t.external_key,
			t.status,
			t.is_active,
			t.final_estimate_value,
			t.created_at,
			t.updated_at,
			COALESCE((
				SELECT GREATEST(
					EXTRACT(EPOCH FROM (MAX(tr.updated_at) - MIN(tr.created_at))),
					0
				)::bigint
				FROM task_rounds AS tr
				WHERE tr.task_id = t.task_id
			), 0)::bigint AS approx_duration_seconds,
			COALESCE((
				SELECT COUNT(*)
				FROM task_rounds AS tr
				WHERE tr.task_id = t.task_id
			), 0)::int AS round_count
		FROM tasks AS t
		WHERE t.room_id = ?
		ORDER BY t.created_at ASC, t.task_id ASC
	`

	if err := r.db.NewRaw(query, roomID).Scan(ctx, &tasks); err != nil {
		return nil, err
	}

	for idx := range tasks {
		tasks[idx].Rounds = make([]historydto.RoomSummaryTaskRound, 0)
	}

	return tasks, nil
}

func (r *historyRepository) getRoomSummaryRounds(
	ctx context.Context,
	roomID string,
) ([]historydto.RoomSummaryTaskRound, error) {
	rounds := make([]historydto.RoomSummaryTaskRound, 0)
	query := `
		SELECT
			tr.task_id,
			tr.round_number,
			tr.status,
			tr.created_at,
			tr.updated_at,
			tr.eligible_participant_ids
		FROM task_rounds AS tr
		JOIN tasks AS t ON t.task_id = tr.task_id
		WHERE t.room_id = ?
		ORDER BY t.created_at ASC, tr.round_number ASC
	`

	if err := r.db.NewRaw(query, roomID).Scan(ctx, &rounds); err != nil {
		return nil, err
	}

	return rounds, nil
}

func (r *historyRepository) getRoomSummaryVotes(
	ctx context.Context,
	roomID string,
) ([]historydto.RoomSummaryVote, error) {
	votes := make([]historydto.RoomSummaryVote, 0)
	query := `
		SELECT
			v.task_id,
			v.round_number,
			v.participant_id,
			rp.user_id,
			rp.guest_name,
			u.email,
			NULLIF(u.display_name, '') AS display_name,
			u.avatar_url,
			v.value,
			v.created_at
		FROM votes AS v
		JOIN tasks AS t ON t.task_id = v.task_id
		JOIN task_rounds AS tr ON tr.task_id = v.task_id AND tr.round_number = v.round_number
		JOIN room_participants AS rp ON rp.room_participants_id = v.participant_id
		LEFT JOIN users AS u ON u.user_id = rp.user_id
		WHERE t.room_id = ?
		  AND tr.status = 'REVEALED'
		ORDER BY t.created_at ASC, v.round_number ASC, v.created_at ASC, v.votes_id ASC
	`

	if err := r.db.NewRaw(query, roomID).Scan(ctx, &votes); err != nil {
		return nil, err
	}

	return votes, nil
}

func buildListMySessionsWhere(userID string, query historydto.MySessionsQuery) (string, []any) {
	whereSQL := `
		(
			r.admin_user_id = ?
			OR EXISTS (
				SELECT 1
				FROM room_participants AS rp_user
				WHERE rp_user.room_id = r.room_id
				  AND rp_user.user_id = ?
			)
		)
	`

	args := []any{
		userID,
		userID,
	}

	if query.Status != historydto.SessionStatusAll {
		whereSQL += `
		AND r.status = ?
	`

		args = append(args, string(query.Status))
	}

	whereSQL += `
		AND (
			? = 'ALL'
			OR (? = 'ADMIN' AND r.admin_user_id = ?)
			OR (
				? = 'PARTICIPANT'
				AND r.admin_user_id <> ?
				AND EXISTS (
					SELECT 1
					FROM room_participants AS rp_participant
					WHERE rp_participant.room_id = r.room_id
					  AND rp_participant.user_id = ?
				)
			)
		)
	`

	args = append(args,
		string(query.Role),
		string(query.Role),
		userID,
		string(query.Role),
		userID,
		userID,
	)

	return whereSQL, args
}

func buildListTeamSessionsWhere(teamID string, query historydto.TeamSessionsQuery) (string, []any) {
	whereSQL := `
		r.team_id = ?
	`

	args := []any{
		teamID,
	}

	if query.Status != historydto.SessionStatusAll {
		whereSQL += `
		AND r.status = ?
	`

		args = append(args, string(query.Status))
	}

	return whereSQL, args
}

func roomTaskRoundKey(taskID string, roundNumber int) string {
	return taskID + "#" + strconv.Itoa(roundNumber)
}
