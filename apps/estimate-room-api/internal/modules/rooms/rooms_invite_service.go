package rooms

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	roomsmodels "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/models"
	roomsrepositories "github.com/master-bogdan/estimate-room-api/internal/modules/rooms/repositories"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

type RoomsInviteService interface {
	GenerateInviteToken(roomID, roomCode string) (string, error)
	Join(roomID, token string, userID *string, guestName *string) (*JoinRoomResult, error)
	ValidateGuestRoomAccess(roomID, guestToken string) (*roomsmodels.RoomParticipantModel, error)
}

type roomsInviteService struct {
	roomsRepo       roomsrepositories.RoomsRepository
	participantRepo roomsrepositories.RoomParticipantRepository
	tokenKey        []byte
	logger          *slog.Logger
}

const (
	guestAccessCookieName = "room_guest_token"
	guestTokenTTL         = 30 * 24 * time.Hour
)

type JoinRoomResult struct {
	Room        *roomsmodels.RoomsModel           `json:"room"`
	Participant *roomsmodels.RoomParticipantModel `json:"participant"`
	GuestToken  string                            `json:"-"`
}

type roomInviteTokenClaims struct {
	RoomID   string `json:"roomId"`
	RoomCode string `json:"roomCode"`
}

type roomGuestTokenClaims struct {
	RoomID        string                          `json:"roomId"`
	RoomCode      string                          `json:"roomCode"`
	ParticipantID string                          `json:"participantId"`
	GuestName     string                          `json:"guestName"`
	Role          roomsmodels.RoomParticipantRole `json:"role"`
	ExpiresAt     time.Time                       `json:"expiresAt"`
}

func NewRoomsInviteService(
	roomsRepo roomsrepositories.RoomsRepository,
	participantRepo roomsrepositories.RoomParticipantRepository,
	tokenKey string,
) RoomsInviteService {
	return &roomsInviteService{
		roomsRepo:       roomsRepo,
		participantRepo: participantRepo,
		tokenKey:        []byte(tokenKey),
		logger:          logger.L().With(slog.String("service", "rooms-invites")),
	}
}

func (s *roomsInviteService) GenerateInviteToken(roomID, roomCode string) (string, error) {
	return s.generateInviteToken(roomID, roomCode)
}

func (s *roomsInviteService) Join(
	roomID, token string,
	userID *string,
	guestName *string,
) (*JoinRoomResult, error) {
	room, err := s.findRoomByInviteToken(roomID, token)
	if err != nil {
		return nil, err
	}

	if room.Status != "ACTIVE" {
		return nil, apperrors.ErrForbidden
	}

	if userID != nil && strings.TrimSpace(*userID) != "" {
		return s.joinRegisteredUser(room, strings.TrimSpace(*userID))
	}

	if guestName == nil {
		return nil, apperrors.ErrBadRequest
	}

	return s.joinGuest(room, *guestName)
}

func (s *roomsInviteService) findRoomByInviteToken(roomID, token string) (*roomsmodels.RoomsModel, error) {
	trimmedToken := strings.TrimSpace(token)
	if trimmedToken == "" {
		return nil, apperrors.ErrBadRequest
	}

	claims, err := s.parseInviteToken(trimmedToken)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	if claims.RoomID != roomID {
		return nil, apperrors.ErrNotFound
	}

	room, err := s.roomsRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}

	if room.Code != claims.RoomCode {
		return nil, apperrors.ErrNotFound
	}

	return room, nil
}

func (s *roomsInviteService) joinRegisteredUser(room *roomsmodels.RoomsModel, userID string) (*JoinRoomResult, error) {
	participant, err := s.participantRepo.FindActiveByUserID(room.RoomID, userID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	if participant == nil {
		participant, err = s.participantRepo.Create(&roomsmodels.RoomParticipantModel{
			RoomParticipantID: uuid.NewString(),
			RoomID:            room.RoomID,
			UserID:            &userID,
			Role:              roomsmodels.RoomParticipantRoleMember,
		})
		if err != nil {
			return nil, err
		}
	}

	fullRoom, err := s.roomsRepo.FindByID(room.RoomID)
	if err != nil {
		return nil, err
	}

	return &JoinRoomResult{
		Room:        fullRoom,
		Participant: participant,
	}, nil
}

func (s *roomsInviteService) joinGuest(room *roomsmodels.RoomsModel, guestName string) (*JoinRoomResult, error) {
	trimmedGuestName := strings.TrimSpace(guestName)
	if trimmedGuestName == "" {
		return nil, apperrors.ErrBadRequest
	}

	participant, err := s.participantRepo.FindActiveByGuestName(room.RoomID, trimmedGuestName)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}
	if participant != nil {
		return nil, apperrors.ErrConflict
	}

	participant, err = s.participantRepo.Create(&roomsmodels.RoomParticipantModel{
		RoomParticipantID: uuid.NewString(),
		RoomID:            room.RoomID,
		GuestName:         &trimmedGuestName,
		Role:              roomsmodels.RoomParticipantRoleGuest,
	})
	if err != nil {
		return nil, err
	}

	guestToken, err := s.generateGuestToken(room, participant)
	if err != nil {
		return nil, err
	}

	fullRoom, err := s.roomsRepo.FindByID(room.RoomID)
	if err != nil {
		return nil, err
	}

	return &JoinRoomResult{
		Room:        fullRoom,
		Participant: participant,
		GuestToken:  guestToken,
	}, nil
}

func (s *roomsInviteService) ValidateGuestRoomAccess(roomID, guestToken string) (*roomsmodels.RoomParticipantModel, error) {
	claims, err := s.parseGuestToken(guestToken)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	if claims.RoomID != roomID || claims.Role != roomsmodels.RoomParticipantRoleGuest {
		return nil, apperrors.ErrForbidden
	}

	if !claims.ExpiresAt.IsZero() && claims.ExpiresAt.Before(time.Now()) {
		return nil, apperrors.ErrUnauthorized
	}

	participant, err := s.participantRepo.FindActiveByID(roomID, claims.ParticipantID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrForbidden
		}
		return nil, err
	}

	if participant.Role != roomsmodels.RoomParticipantRoleGuest {
		return nil, apperrors.ErrForbidden
	}

	return participant, nil
}

func (s *roomsInviteService) generateInviteToken(roomID, roomCode string) (string, error) {
	return utils.GenerateToken(s.tokenKey, roomInviteTokenClaims{
		RoomID:   roomID,
		RoomCode: roomCode,
	})
}

func (s *roomsInviteService) parseInviteToken(token string) (*roomInviteTokenClaims, error) {
	return utils.ParseToken[roomInviteTokenClaims](s.tokenKey, token)
}

func (s *roomsInviteService) generateGuestToken(
	room *roomsmodels.RoomsModel,
	participant *roomsmodels.RoomParticipantModel,
) (string, error) {
	guestName := ""
	if participant.GuestName != nil {
		guestName = *participant.GuestName
	}

	return utils.GenerateToken(s.tokenKey, roomGuestTokenClaims{
		RoomID:        room.RoomID,
		RoomCode:      room.Code,
		ParticipantID: participant.RoomParticipantID,
		GuestName:     guestName,
		Role:          roomsmodels.RoomParticipantRoleGuest,
		ExpiresAt:     time.Now().Add(guestTokenTTL),
	})
}

func (s *roomsInviteService) parseGuestToken(token string) (*roomGuestTokenClaims, error) {
	return utils.ParseToken[roomGuestTokenClaims](s.tokenKey, token)
}
