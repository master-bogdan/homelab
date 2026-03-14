package rooms

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/modules/ws"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type roomsGateway struct {
	wsService       *ws.Service
	roomsRepo       roomsrepositories.RoomsRepository
	participantRepo roomsrepositories.RoomParticipantRepository
	taskRepo        roomsrepositories.RoomTaskRepository
	voteRepo        roomsrepositories.RoomVoteRepository
	roundRepo       roomsrepositories.RoomTaskRoundRepository
	voteService     RoomsVoteService
	expiryService   RoomsExpiryService
}

func NewRoomsGateway(
	wsService *ws.Service,
	roomsRepo roomsrepositories.RoomsRepository,
	participantRepo roomsrepositories.RoomParticipantRepository,
	taskRepo roomsrepositories.RoomTaskRepository,
	voteRepo roomsrepositories.RoomVoteRepository,
	roundRepo roomsrepositories.RoomTaskRoundRepository,
	voteService RoomsVoteService,
	expiryService RoomsExpiryService,
) *roomsGateway {
	return &roomsGateway{
		wsService:       wsService,
		roomsRepo:       roomsRepo,
		participantRepo: participantRepo,
		taskRepo:        taskRepo,
		voteRepo:        voteRepo,
		roundRepo:       roundRepo,
		voteService:     voteService,
		expiryService:   expiryService,
	}
}

const (
	RoomsJoin           = "ROOMS_JOIN"
	RoomsTaskSetCurrent = "ROOMS_TASK_SET_CURRENT"
	RoomsVoteCast       = "ROOMS_VOTE_CAST"
	RoomsVoteReveal     = "ROOMS_VOTE_REVEAL"
	RoomsRoundNext      = "ROOMS_ROUND_NEXT"
	RoomsTaskFinalize   = "ROOMS_TASK_FINALIZE"

	RoomsParticipantJoined  = "ROOMS_PARTICIPANT_JOINED"
	RoomsParticipantLeft    = "ROOMS_PARTICIPANT_LEFT"
	RoomsTaskCurrentChanged = "ROOMS_TASK_CURRENT_CHANGED"
	RoomsVoteStatusChanged  = "ROOMS_VOTE_STATUS_CHANGED"
	RoomsVotesAllCast       = "ROOMS_VOTES_ALL_CAST"
	RoomsVotesRevealed      = "ROOMS_VOTES_REVEALED"
	RoomsRoundChanged       = "ROOMS_ROUND_CHANGED"
	RoomsTaskFinalized      = "ROOMS_TASK_FINALIZED"
	RoomsExpired            = "ROOMS_EXPIRED"
	RoomsSnapshot           = "ROOMS_SNAPSHOT"
)

type roomJoinPayload struct {
	RoomID string `json:"roomId"`
}

type roomPresencePayload struct {
	ParticipantID string                          `json:"participantId,omitempty"`
	UserID        *string                         `json:"userId,omitempty"`
	GuestName     *string                         `json:"guestName,omitempty"`
	Role          roomsmodels.RoomParticipantRole `json:"role,omitempty"`
}

type roomSetCurrentTaskPayload struct {
	TaskID string `json:"taskId"`
}

type roomCurrentTaskChangedPayload struct {
	CurrentTaskID          string   `json:"currentTaskId"`
	PreviousTaskID         *string  `json:"previousTaskId,omitempty"`
	RoundNumber            int      `json:"roundNumber"`
	RoundStatus            string   `json:"roundStatus"`
	EligibleParticipantIDs []string `json:"eligibleParticipantIds"`
}

type roomVoteCastPayload struct {
	Value string `json:"value"`
}

type roomVoteStatusChangedPayload struct {
	TaskID        string `json:"taskId"`
	ParticipantID string `json:"participantId"`
	RoundNumber   int    `json:"roundNumber"`
	Voted         bool   `json:"voted"`
}

type roomVotesAllCastPayload struct {
	TaskID                 string   `json:"taskId"`
	RoundNumber            int      `json:"roundNumber"`
	EligibleParticipantIDs []string `json:"eligibleParticipantIds"`
	VotedParticipantIDs    []string `json:"votedParticipantIds"`
}

type roomRevealedVote struct {
	ParticipantID string `json:"participantId"`
	Value         string `json:"value"`
}

type roomVoteSummary struct {
	TotalVotes int            `json:"totalVotes"`
	Counts     map[string]int `json:"counts"`
}

type roomVotesRevealedPayload struct {
	TaskID      string             `json:"taskId"`
	RoundNumber int                `json:"roundNumber"`
	RoundStatus string             `json:"roundStatus"`
	AllVoted    bool               `json:"allVoted"`
	Votes       []roomRevealedVote `json:"votes"`
	Summary     roomVoteSummary    `json:"summary"`
}

type roomRoundChangedPayload struct {
	TaskID                 string   `json:"taskId"`
	RoundNumber            int      `json:"roundNumber"`
	RoundStatus            string   `json:"roundStatus"`
	EligibleParticipantIDs []string `json:"eligibleParticipantIds"`
}

type roomTaskFinalizePayload struct {
	Value string `json:"value"`
}

type roomTaskFinalizedPayload struct {
	TaskID             string `json:"taskId"`
	FinalEstimateValue string `json:"finalEstimateValue"`
	Status             string `json:"status"`
}

type roomSnapshotRoom struct {
	RoomID      string               `json:"roomId"`
	Code        string               `json:"code"`
	Name        string               `json:"name"`
	Status      string               `json:"status"`
	AdminUserID string               `json:"adminUserId"`
	Deck        roomsmodels.RoomDeck `json:"deck"`
}

type roomSnapshotParticipant struct {
	ParticipantID string                          `json:"participantId"`
	UserID        *string                         `json:"userId,omitempty"`
	GuestName     *string                         `json:"guestName,omitempty"`
	Role          roomsmodels.RoomParticipantRole `json:"role"`
	Online        bool                            `json:"online"`
}

type roomSnapshotTask struct {
	TaskID             string  `json:"taskId"`
	Title              string  `json:"title"`
	Description        *string `json:"description,omitempty"`
	ExternalKey        *string `json:"externalKey,omitempty"`
	Status             string  `json:"status"`
	FinalEstimateValue *string `json:"finalEstimateValue,omitempty"`
}

type roomSnapshotPayload struct {
	Room                   roomSnapshotRoom          `json:"room"`
	Participants           []roomSnapshotParticipant `json:"participants"`
	Tasks                  []roomSnapshotTask        `json:"tasks"`
	CurrentTaskID          *string                   `json:"currentTaskId,omitempty"`
	CurrentRoundNumber     int                       `json:"currentRoundNumber"`
	RoundStatus            string                    `json:"roundStatus"`
	EligibleParticipantIDs []string                  `json:"eligibleParticipantIds"`
	VotedParticipantIDs    []string                  `json:"votedParticipantIds"`
	RevealedVotes          []roomRevealedVote        `json:"revealedVotes,omitempty"`
	Summary                *roomVoteSummary          `json:"summary,omitempty"`
}

func (g *roomsGateway) handleRoomJoin(client ws.ClientInfo, event ws.Event) {
	roomID := resolveRoomID(event)
	if roomID == "" {
		logger.L().Warn("room join ignored: missing room id", "user_id", client.UserID, "conn_id", client.ConnID)
		return
	}

	participant, err := g.resolveParticipant(client, roomID)
	if err != nil {
		logJoinDenied(client, roomID, err)
		return
	}

	if err := g.wsService.SetParticipantID(client.ConnID, participant.RoomParticipantID); err != nil {
		logger.L().Error("failed to bind ws participant", "err", err, "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	joinResult, err := g.wsService.JoinRoom(client.ConnID, roomID)
	if err != nil {
		logger.L().Error("room join failed", "err", err, "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	g.expiryService.TouchActivity(roomID)

	if joinResult.Joined {
		if err := g.broadcastPresence(roomID, RoomsParticipantJoined, roomPresencePayload{
			ParticipantID: participant.RoomParticipantID,
			UserID:        participant.UserID,
			GuestName:     participant.GuestName,
			Role:          participant.Role,
		}); err != nil {
			logger.L().Error("failed to broadcast participant joined", "err", err, "room_id", roomID)
		}
	}

	snapshot, err := g.buildSnapshot(roomID)
	if err != nil {
		logger.L().Error("failed to build room snapshot", "err", err, "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	if err := g.sendSnapshot(client.ConnID, roomID, snapshot); err != nil {
		logger.L().Error("failed to send room snapshot", "err", err, "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	logger.L().Info("room join accepted", "room_id", roomID, "user_id", client.UserID, "conn_id", client.ConnID)
}

func (g *roomsGateway) handleTaskSetCurrent(client ws.ClientInfo, event ws.Event) {
	roomID := strings.TrimSpace(event.RoomID)
	if roomID == "" {
		logger.L().Warn("task set current ignored: missing room id", "user_id", client.UserID, "conn_id", client.ConnID)
		return
	}

	payload := roomSetCurrentTaskPayload{}
	if len(event.Payload) > 0 {
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			logger.L().Warn("task set current ignored: invalid payload", "err", err, "room_id", roomID, "conn_id", client.ConnID)
			return
		}
	}

	taskID := strings.TrimSpace(payload.TaskID)
	if taskID == "" {
		logger.L().Warn("task set current ignored: missing task id", "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	participant, err := g.resolveParticipant(client, roomID)
	if err != nil {
		logJoinDenied(client, roomID, err)
		return
	}
	if participant.Role != roomsmodels.RoomParticipantRoleAdmin {
		logger.L().Warn("task set current denied: admin only", "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	eligibleParticipantIDs, err := g.currentEligibleParticipantIDs(roomID)
	if err != nil {
		logger.L().Error("task set current failed: eligible participants lookup failed", "room_id", roomID, "task_id", taskID, "err", err)
		return
	}

	currentTask, previousTask, roundState, err := g.voteService.SetCurrentTask(roomID, taskID, client.UserID, eligibleParticipantIDs)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, apperrors.ErrForbidden) || errors.Is(err, apperrors.ErrBadRequest) {
			logger.L().Warn("task set current denied", "room_id", roomID, "task_id", taskID, "reason", err.Error())
			return
		}
		logger.L().Error("task set current failed", "room_id", roomID, "task_id", taskID, "err", err)
		return
	}

	var previousTaskID *string
	if previousTask != nil && strings.TrimSpace(previousTask.TaskID) != "" {
		id := previousTask.TaskID
		previousTaskID = &id
	}

	if err := g.broadcastCurrentTaskChanged(roomID, roomCurrentTaskChangedPayload{
		CurrentTaskID:          currentTask.TaskID,
		PreviousTaskID:         previousTaskID,
		RoundNumber:            roundState.RoundNumber,
		RoundStatus:            string(roundState.Status),
		EligibleParticipantIDs: append([]string(nil), roundState.EligibleParticipantIDs...),
	}); err != nil {
		logger.L().Error("failed to broadcast current task changed", "room_id", roomID, "task_id", taskID, "err", err)
		return
	}

	logger.L().Info("task current changed", "room_id", roomID, "task_id", currentTask.TaskID, "conn_id", client.ConnID)
}

func (g *roomsGateway) handleVoteCast(client ws.ClientInfo, event ws.Event) {
	roomID := strings.TrimSpace(event.RoomID)
	if roomID == "" {
		logger.L().Warn("vote cast ignored: missing room id", "user_id", client.UserID, "conn_id", client.ConnID)
		return
	}

	payload := roomVoteCastPayload{}
	if len(event.Payload) > 0 {
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			logger.L().Warn("vote cast ignored: invalid payload", "err", err, "room_id", roomID, "conn_id", client.ConnID)
			return
		}
	}

	voteValue := strings.TrimSpace(payload.Value)
	if voteValue == "" {
		logger.L().Warn("vote cast ignored: empty vote value", "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	participant, err := g.resolveParticipant(client, roomID)
	if err != nil {
		logJoinDenied(client, roomID, err)
		return
	}

	result, err := g.voteService.CastVote(roomID, participant, voteValue)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrForbidden), errors.Is(err, apperrors.ErrBadRequest):
			logger.L().Warn("vote cast denied", "room_id", roomID, "conn_id", client.ConnID, "reason", err.Error())
		default:
			logger.L().Error("vote cast failed", "room_id", roomID, "conn_id", client.ConnID, "err", err)
		}
		return
	}

	if err := g.broadcastVoteStatusChanged(roomID, roomVoteStatusChangedPayload{
		TaskID:        result.Task.TaskID,
		ParticipantID: participant.RoomParticipantID,
		RoundNumber:   result.Round.RoundNumber,
		Voted:         true,
	}); err != nil {
		logger.L().Error("failed to broadcast vote status changed", "room_id", roomID, "task_id", result.Task.TaskID, "err", err)
		return
	}

	if result.AllVotesCast {
		if err := g.broadcastVotesAllCast(roomID, roomVotesAllCastPayload{
			TaskID:                 result.Task.TaskID,
			RoundNumber:            result.Round.RoundNumber,
			EligibleParticipantIDs: append([]string(nil), result.EligibleParticipantIDs...),
			VotedParticipantIDs:    append([]string(nil), result.VotedParticipantIDs...),
		}); err != nil {
			logger.L().Error("failed to broadcast votes all cast", "room_id", roomID, "task_id", result.Task.TaskID, "err", err)
		}
	}
}

func (g *roomsGateway) handleVoteReveal(client ws.ClientInfo, event ws.Event) {
	roomID := strings.TrimSpace(event.RoomID)
	if roomID == "" {
		logger.L().Warn("vote reveal ignored: missing room id", "user_id", client.UserID, "conn_id", client.ConnID)
		return
	}

	participant, err := g.resolveParticipant(client, roomID)
	if err != nil {
		logJoinDenied(client, roomID, err)
		return
	}
	if participant.Role != roomsmodels.RoomParticipantRoleAdmin {
		logger.L().Warn("vote reveal denied: admin only", "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	result, err := g.voteService.RevealCurrentRound(roomID, client.UserID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrForbidden), errors.Is(err, apperrors.ErrBadRequest):
			logger.L().Warn("vote reveal denied", "room_id", roomID, "conn_id", client.ConnID, "reason", err.Error())
		default:
			logger.L().Error("vote reveal failed", "room_id", roomID, "conn_id", client.ConnID, "err", err)
		}
		return
	}

	revealedVotes := mapVotes(result.Votes)
	summary := buildVoteSummary(result.Votes)
	payload := roomVotesRevealedPayload{
		TaskID:      result.Task.TaskID,
		RoundNumber: result.Round.RoundNumber,
		RoundStatus: string(result.Round.Status),
		AllVoted:    result.AllVoted,
		Votes:       revealedVotes,
		Summary:     summary,
	}

	if err := g.broadcastVotesRevealed(roomID, payload); err != nil {
		logger.L().Error("failed to broadcast votes revealed", "room_id", roomID, "task_id", result.Task.TaskID, "round", result.Round.RoundNumber, "err", err)
		return
	}
}

func (g *roomsGateway) handleRoundNext(client ws.ClientInfo, event ws.Event) {
	roomID := strings.TrimSpace(event.RoomID)
	if roomID == "" {
		logger.L().Warn("round next ignored: missing room id", "user_id", client.UserID, "conn_id", client.ConnID)
		return
	}

	participant, err := g.resolveParticipant(client, roomID)
	if err != nil {
		logJoinDenied(client, roomID, err)
		return
	}
	if participant.Role != roomsmodels.RoomParticipantRoleAdmin {
		logger.L().Warn("round next denied: admin only", "room_id", roomID, "conn_id", client.ConnID)
		return
	}

	eligibleParticipantIDs, err := g.currentEligibleParticipantIDs(roomID)
	if err != nil {
		logger.L().Error("round next failed: eligible participants lookup failed", "room_id", roomID, "conn_id", client.ConnID, "err", err)
		return
	}

	currentTask, nextRound, err := g.voteService.StartNextRound(roomID, client.UserID, eligibleParticipantIDs)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrForbidden), errors.Is(err, apperrors.ErrBadRequest):
			logger.L().Warn("round next denied", "room_id", roomID, "conn_id", client.ConnID, "reason", err.Error())
		default:
			logger.L().Error("round next failed", "room_id", roomID, "conn_id", client.ConnID, "err", err)
		}
		return
	}

	if err := g.broadcastRoundChanged(roomID, roomRoundChangedPayload{
		TaskID:                 currentTask.TaskID,
		RoundNumber:            nextRound.RoundNumber,
		RoundStatus:            string(nextRound.Status),
		EligibleParticipantIDs: append([]string(nil), nextRound.EligibleParticipantIDs...),
	}); err != nil {
		logger.L().Error("failed to broadcast round changed", "room_id", roomID, "task_id", currentTask.TaskID, "err", err)
	}
}

func (g *roomsGateway) handleTaskFinalize(client ws.ClientInfo, event ws.Event) {
	roomID := strings.TrimSpace(event.RoomID)
	if roomID == "" {
		logger.L().Warn("task finalize ignored: missing room id", "user_id", client.UserID, "conn_id", client.ConnID)
		return
	}

	payload := roomTaskFinalizePayload{}
	if len(event.Payload) > 0 {
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			logger.L().Warn("task finalize ignored: invalid payload", "err", err, "room_id", roomID, "conn_id", client.ConnID)
			return
		}
	}

	updatedTask, err := g.voteService.FinalizeCurrentTask(roomID, client.UserID, payload.Value)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrNotFound), errors.Is(err, apperrors.ErrForbidden), errors.Is(err, apperrors.ErrBadRequest):
			logger.L().Warn("task finalize denied", "room_id", roomID, "conn_id", client.ConnID, "reason", err.Error())
		default:
			logger.L().Error("task finalize failed", "room_id", roomID, "conn_id", client.ConnID, "err", err)
		}
		return
	}

	if err := g.broadcastTaskFinalized(roomID, roomTaskFinalizedPayload{
		TaskID:             updatedTask.TaskID,
		FinalEstimateValue: strings.TrimSpace(*updatedTask.FinalEstimateValue),
		Status:             updatedTask.Status,
	}); err != nil {
		logger.L().Error("failed to broadcast task finalized", "room_id", roomID, "task_id", updatedTask.TaskID, "err", err)
	}
}

func (g *roomsGateway) handleDisconnect(info ws.DisconnectInfo) {
	roomID := strings.TrimSpace(info.RoomID)
	if roomID == "" || !info.PresenceLeft {
		return
	}

	participantID := strings.TrimSpace(info.Client.ParticipantID)
	payload := roomPresencePayload{
		ParticipantID: participantID,
	}

	if info.Client.UserID != "" {
		userID := info.Client.UserID
		payload.UserID = &userID
	}

	if err := g.broadcastPresence(roomID, RoomsParticipantLeft, payload); err != nil {
		logger.L().Error("failed to broadcast participant left", "err", err, "room_id", roomID)
	}

	g.expiryService.TouchActivity(roomID)
}

func (g *roomsGateway) resolveParticipant(client ws.ClientInfo, roomID string) (*roomsmodels.RoomParticipantModel, error) {
	switch client.IdentityType {
	case ws.IdentityTypeUser:
		userID := strings.TrimSpace(client.UserID)
		if userID == "" {
			return nil, apperrors.ErrUnauthorized
		}
		return g.participantRepo.FindActiveByUserID(roomID, userID)
	case ws.IdentityTypeGuest:
		participantID := strings.TrimSpace(client.ParticipantID)
		if participantID == "" {
			return nil, apperrors.ErrUnauthorized
		}
		participant, err := g.participantRepo.FindActiveByID(roomID, participantID)
		if err != nil {
			return nil, err
		}
		if participant.Role != roomsmodels.RoomParticipantRoleGuest {
			return nil, apperrors.ErrForbidden
		}
		return participant, nil
	default:
		return nil, apperrors.ErrUnauthorized
	}
}

func (g *roomsGateway) buildSnapshot(roomID string) (*roomSnapshotPayload, error) {
	room, err := g.roomsRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	onlineIDs := g.wsService.GetRoomOnlineParticipantIDs(roomID)
	onlineSet := make(map[string]struct{}, len(onlineIDs))
	for _, id := range onlineIDs {
		onlineSet[id] = struct{}{}
	}

	participants := make([]roomSnapshotParticipant, 0, len(room.Participants))
	for _, participant := range room.Participants {
		_, online := onlineSet[participant.RoomParticipantID]
		participants = append(participants, roomSnapshotParticipant{
			ParticipantID: participant.RoomParticipantID,
			UserID:        participant.UserID,
			GuestName:     participant.GuestName,
			Role:          participant.Role,
			Online:        online,
		})
	}

	tasks := make([]roomSnapshotTask, 0, len(room.Tasks))
	for _, task := range room.Tasks {
		tasks = append(tasks, roomSnapshotTask{
			TaskID:             task.TaskID,
			Title:              task.Title,
			Description:        task.Description,
			ExternalKey:        task.ExternalKey,
			Status:             task.Status,
			FinalEstimateValue: task.FinalEstimateValue,
		})
	}

	snapshot := &roomSnapshotPayload{
		Room: roomSnapshotRoom{
			RoomID:      room.RoomID,
			Code:        room.Code,
			Name:        room.Name,
			Status:      room.Status,
			AdminUserID: room.AdminUserID,
			Deck:        room.Deck,
		},
		Participants:           participants,
		Tasks:                  tasks,
		CurrentRoundNumber:     1,
		RoundStatus:            "",
		EligibleParticipantIDs: make([]string, 0),
		VotedParticipantIDs:    make([]string, 0),
	}

	currentTask, err := g.taskRepo.FindCurrentVotingTask(roomID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return snapshot, nil
		}
		return nil, err
	}

	currentTaskID := currentTask.TaskID
	snapshot.CurrentTaskID = &currentTaskID

	currentRound, err := g.roundRepo.GetOrCreateCurrent(currentTask.TaskID, nil)
	if err != nil {
		return nil, err
	}
	snapshot.CurrentRoundNumber = currentRound.RoundNumber
	snapshot.RoundStatus = string(currentRound.Status)
	snapshot.EligibleParticipantIDs = append([]string(nil), currentRound.EligibleParticipantIDs...)

	votes, err := g.voteRepo.ListByTaskAndRound(currentTask.TaskID, currentRound.RoundNumber)
	if err != nil {
		return nil, err
	}

	snapshot.VotedParticipantIDs = filterParticipantIDs(uniqueSortedParticipantIDs(votes), currentRound.EligibleParticipantIDs)
	if currentRound.Status == roomsmodels.RoomTaskRoundStatusRevealed {
		revealedVotes := mapVotes(votes)
		summary := buildVoteSummary(votes)
		snapshot.RevealedVotes = revealedVotes
		snapshot.Summary = &summary
	}

	return snapshot, nil
}

func (g *roomsGateway) currentEligibleParticipantIDs(roomID string) ([]string, error) {
	onlineParticipantIDs := g.wsService.GetRoomOnlineParticipantIDs(roomID)
	if len(onlineParticipantIDs) == 0 {
		return []string{}, nil
	}

	participants, err := g.participantRepo.ListActiveByRoom(roomID)
	if err != nil {
		return nil, err
	}

	eligible := make([]string, 0, len(participants))
	for _, participant := range participants {
		if participant == nil || !isVotingParticipantRole(participant.Role) {
			continue
		}
		if containsParticipantID(onlineParticipantIDs, participant.RoomParticipantID) {
			eligible = append(eligible, participant.RoomParticipantID)
		}
	}

	return normalizeParticipantIDs(eligible), nil
}

func (g *roomsGateway) sendSnapshot(connID, roomID string, snapshot *roomSnapshotPayload) error {
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	return g.wsService.SendToConnection(connID, ws.Event{
		Type:    RoomsSnapshot,
		RoomID:  roomID,
		Payload: data,
	})
}

func resolveRoomID(event ws.Event) string {
	roomID := strings.TrimSpace(event.RoomID)
	if roomID != "" {
		return roomID
	}

	if len(event.Payload) == 0 {
		return ""
	}

	payload := roomJoinPayload{}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return ""
	}
	return strings.TrimSpace(payload.RoomID)
}

func (g *roomsGateway) broadcastPresence(roomID, eventType string, payload roomPresencePayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.wsService.Broadcast(ws.Event{
		Type:    eventType,
		RoomID:  roomID,
		Payload: data,
	})
}

func (g *roomsGateway) broadcastCurrentTaskChanged(roomID string, payload roomCurrentTaskChangedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.wsService.Broadcast(ws.Event{
		Type:    RoomsTaskCurrentChanged,
		RoomID:  roomID,
		Payload: data,
	})
}

func (g *roomsGateway) broadcastVoteStatusChanged(roomID string, payload roomVoteStatusChangedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.wsService.Broadcast(ws.Event{
		Type:    RoomsVoteStatusChanged,
		RoomID:  roomID,
		Payload: data,
	})
}

func (g *roomsGateway) broadcastVotesAllCast(roomID string, payload roomVotesAllCastPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.wsService.Broadcast(ws.Event{
		Type:    RoomsVotesAllCast,
		RoomID:  roomID,
		Payload: data,
	})
}

func (g *roomsGateway) broadcastVotesRevealed(roomID string, payload roomVotesRevealedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.wsService.Broadcast(ws.Event{
		Type:    RoomsVotesRevealed,
		RoomID:  roomID,
		Payload: data,
	})
}

func (g *roomsGateway) broadcastRoundChanged(roomID string, payload roomRoundChangedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.wsService.Broadcast(ws.Event{
		Type:    RoomsRoundChanged,
		RoomID:  roomID,
		Payload: data,
	})
}

func (g *roomsGateway) broadcastTaskFinalized(roomID string, payload roomTaskFinalizedPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return g.wsService.Broadcast(ws.Event{
		Type:    RoomsTaskFinalized,
		RoomID:  roomID,
		Payload: data,
	})
}

func mapVotes(votes []*roomsmodels.RoomVoteModel) []roomRevealedVote {
	revealed := make([]roomRevealedVote, 0, len(votes))
	for _, vote := range votes {
		revealed = append(revealed, roomRevealedVote{
			ParticipantID: vote.ParticipantID,
			Value:         vote.Value,
		})
	}
	return revealed
}

func buildVoteSummary(votes []*roomsmodels.RoomVoteModel) roomVoteSummary {
	counts := make(map[string]int, len(votes))
	for _, vote := range votes {
		counts[vote.Value]++
	}

	return roomVoteSummary{
		TotalVotes: len(votes),
		Counts:     counts,
	}
}

func uniqueSortedParticipantIDs(votes []*roomsmodels.RoomVoteModel) []string {
	set := make(map[string]struct{}, len(votes))
	for _, vote := range votes {
		participantID := strings.TrimSpace(vote.ParticipantID)
		if participantID != "" {
			set[participantID] = struct{}{}
		}
	}

	ids := make([]string, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func isDeckValueAllowed(values []string, value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}

	for _, deckValue := range values {
		if strings.TrimSpace(deckValue) == trimmed {
			return true
		}
	}

	return false
}

func logJoinDenied(client ws.ClientInfo, roomID string, err error) {
	switch {
	case errors.Is(err, apperrors.ErrForbidden), errors.Is(err, apperrors.ErrUnauthorized), errors.Is(err, apperrors.ErrNotFound):
		logger.L().Warn("room join denied", "room_id", roomID, "user_id", client.UserID, "conn_id", client.ConnID, "reason", err.Error())
	default:
		logger.L().Error("room join failed", "room_id", roomID, "user_id", client.UserID, "conn_id", client.ConnID, "err", err)
	}
}
