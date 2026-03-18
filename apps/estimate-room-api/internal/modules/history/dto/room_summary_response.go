package historydto

import "time"

type RoomSummaryResponse struct {
	Overview     RoomSummaryOverview      `json:"overview"`
	Participants []RoomSummaryParticipant `json:"participants"`
	Tasks        []RoomSummaryTask        `json:"tasks"`
}

type RoomSummaryOverview struct {
	RoomID                string              `json:"roomId" bun:"room_id"`
	TeamID                *string             `json:"teamId,omitempty" bun:"team_id"`
	Name                  string              `json:"name" bun:"name"`
	Status                string              `json:"status" bun:"status"`
	CreatedAt             time.Time           `json:"createdAt" bun:"created_at"`
	FinishedAt            *time.Time          `json:"finishedAt,omitempty" bun:"finished_at"`
	LastActivityAt        time.Time           `json:"lastActivityAt" bun:"last_activity_at"`
	ApproxDurationSeconds int64               `json:"approxDurationSeconds" bun:"approx_duration_seconds"`
	ParticipantsCount     int                 `json:"participantsCount" bun:"participants_count"`
	EstimatedTasksCount   int                 `json:"estimatedTasksCount" bun:"estimated_tasks_count"`
	TasksCount            int                 `json:"tasksCount" bun:"tasks_count"`
	RoundCount            int                 `json:"roundCount" bun:"round_count"`
	AdminUser             RoomSummaryUserRef  `json:"adminUser"`
}

type RoomSummaryUserRef struct {
	UserID      string  `json:"userId" bun:"admin_user_id"`
	Email       *string `json:"email,omitempty" bun:"admin_email"`
	DisplayName string  `json:"displayName" bun:"admin_display_name"`
	AvatarURL   *string `json:"avatarUrl,omitempty" bun:"admin_avatar_url"`
}

type RoomSummaryParticipant struct {
	ParticipantID            string     `json:"participantId" bun:"participant_id"`
	UserID                   *string    `json:"userId,omitempty" bun:"user_id"`
	GuestName                *string    `json:"guestName,omitempty" bun:"guest_name"`
	Email                    *string    `json:"email,omitempty" bun:"email"`
	DisplayName              *string    `json:"displayName,omitempty" bun:"display_name"`
	AvatarURL                *string    `json:"avatarUrl,omitempty" bun:"avatar_url"`
	Role                     string     `json:"role" bun:"role"`
	JoinedAt                 time.Time  `json:"joinedAt" bun:"joined_at"`
	LeftAt                   *time.Time `json:"leftAt,omitempty" bun:"left_at"`
	VotesCastCount           int        `json:"votesCastCount" bun:"votes_cast_count"`
	EstimatedTasksVotedCount int        `json:"estimatedTasksVotedCount" bun:"estimated_tasks_voted_count"`
}

type RoomSummaryTask struct {
	TaskID                 string                 `json:"taskId" bun:"task_id"`
	Title                  string                 `json:"title" bun:"title"`
	Description            *string                `json:"description,omitempty" bun:"description"`
	ExternalKey            *string                `json:"externalKey,omitempty" bun:"external_key"`
	Status                 string                 `json:"status" bun:"status"`
	IsActive               bool                   `json:"isActive" bun:"is_active"`
	FinalEstimateValue     *string                `json:"finalEstimateValue,omitempty" bun:"final_estimate_value"`
	CreatedAt              time.Time              `json:"createdAt" bun:"created_at"`
	UpdatedAt              time.Time              `json:"updatedAt" bun:"updated_at"`
	ApproxDurationSeconds  int64                  `json:"approxDurationSeconds" bun:"approx_duration_seconds"`
	RoundCount             int                    `json:"roundCount" bun:"round_count"`
	Rounds                 []RoomSummaryTaskRound `json:"rounds"`
}

type RoomSummaryTaskRound struct {
	TaskID                 string            `json:"-" bun:"task_id"`
	RoundNumber            int               `json:"roundNumber" bun:"round_number"`
	Status                 string            `json:"status" bun:"status"`
	CreatedAt              time.Time         `json:"createdAt" bun:"created_at"`
	UpdatedAt              time.Time         `json:"updatedAt" bun:"updated_at"`
	EligibleParticipantIDs []string          `json:"eligibleParticipantIds" bun:"eligible_participant_ids"`
	Votes                  []RoomSummaryVote `json:"votes"`
}

type RoomSummaryVote struct {
	TaskID        string     `json:"-" bun:"task_id"`
	RoundNumber   int        `json:"-" bun:"round_number"`
	ParticipantID string     `json:"participantId" bun:"participant_id"`
	UserID        *string    `json:"userId,omitempty" bun:"user_id"`
	GuestName     *string    `json:"guestName,omitempty" bun:"guest_name"`
	Email         *string    `json:"email,omitempty" bun:"email"`
	DisplayName   *string    `json:"displayName,omitempty" bun:"display_name"`
	AvatarURL     *string    `json:"avatarUrl,omitempty" bun:"avatar_url"`
	Value         string     `json:"value" bun:"value"`
	CreatedAt     time.Time  `json:"createdAt" bun:"created_at"`
}
