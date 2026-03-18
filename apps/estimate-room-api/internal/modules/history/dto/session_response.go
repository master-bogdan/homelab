package historydto

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type SessionStatusFilter string
type SessionRoleFilter string

const (
	SessionStatusAll      SessionStatusFilter = "ALL"
	SessionStatusActive   SessionStatusFilter = "ACTIVE"
	SessionStatusFinished SessionStatusFilter = "FINISHED"
	SessionStatusExpired  SessionStatusFilter = "EXPIRED"

	SessionRoleAll         SessionRoleFilter = "ALL"
	SessionRoleAdmin       SessionRoleFilter = "ADMIN"
	SessionRoleParticipant SessionRoleFilter = "PARTICIPANT"
)

type MySessionsQuery struct {
	PaginationQuery
	Status SessionStatusFilter `json:"status"`
	Role   SessionRoleFilter   `json:"role"`
}

type TeamSessionsQuery struct {
	PaginationQuery
	Status SessionStatusFilter `json:"status"`
}

type SessionListItem struct {
	RoomID                string     `json:"roomId" bun:"room_id"`
	TeamID                *string    `json:"teamId,omitempty" bun:"team_id"`
	Name                  string     `json:"name" bun:"name"`
	Status                string     `json:"status" bun:"status"`
	Role                  string     `json:"role" bun:"role"`
	CreatedAt             time.Time  `json:"createdAt" bun:"created_at"`
	FinishedAt            *time.Time `json:"finishedAt,omitempty" bun:"finished_at"`
	LastActivityAt        time.Time  `json:"lastActivityAt" bun:"last_activity_at"`
	ApproxDurationSeconds int64      `json:"approxDurationSeconds" bun:"approx_duration_seconds"`
	ParticipantsCount     int        `json:"participantsCount" bun:"participants_count"`
	TasksCount            int        `json:"tasksCount" bun:"tasks_count"`
	EstimatedTasksCount   int        `json:"estimatedTasksCount" bun:"estimated_tasks_count"`
}

type SessionListResponse struct {
	Items    []SessionListItem `json:"items"`
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
	Total    int               `json:"total"`
}

func ParseMySessionsQuery(values url.Values) (MySessionsQuery, error) {
	pagination, err := ParsePaginationQuery(values)
	if err != nil {
		return MySessionsQuery{}, err
	}

	query := MySessionsQuery{
		PaginationQuery: pagination,
		Status:          SessionStatusAll,
		Role:            SessionRoleAll,
	}

	status, err := parseSessionStatusFilter(values)
	if err != nil {
		return MySessionsQuery{}, err
	}
	query.Status = status

	if rawRole := strings.ToUpper(strings.TrimSpace(values.Get("role"))); rawRole != "" {
		switch SessionRoleFilter(rawRole) {
		case SessionRoleAll, SessionRoleAdmin, SessionRoleParticipant:
			query.Role = SessionRoleFilter(rawRole)
		default:
			return MySessionsQuery{}, fmt.Errorf("role must be one of ALL, ADMIN, PARTICIPANT")
		}
	}

	return query, nil
}

func ParseTeamSessionsQuery(values url.Values) (TeamSessionsQuery, error) {
	pagination, err := ParsePaginationQuery(values)
	if err != nil {
		return TeamSessionsQuery{}, err
	}

	status, err := parseSessionStatusFilter(values)
	if err != nil {
		return TeamSessionsQuery{}, err
	}

	return TeamSessionsQuery{
		PaginationQuery: pagination,
		Status:          status,
	}, nil
}

func parseSessionStatusFilter(values url.Values) (SessionStatusFilter, error) {
	rawStatus := strings.ToUpper(strings.TrimSpace(values.Get("status")))
	if rawStatus == "" {
		return SessionStatusAll, nil
	}

	switch SessionStatusFilter(rawStatus) {
	case SessionStatusAll, SessionStatusActive, SessionStatusFinished, SessionStatusExpired:
		return SessionStatusFilter(rawStatus), nil
	default:
		return "", fmt.Errorf("status must be one of ALL, ACTIVE, FINISHED, EXPIRED")
	}
}
