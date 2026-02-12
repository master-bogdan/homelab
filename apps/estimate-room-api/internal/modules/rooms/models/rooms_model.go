package roomsmodels

import "time"

type RoomsModel struct {
	RoomID            string
	Code              string
	Name              string
	AdminUserID       string
	TeamID            *string
	DeckID            DeckID
	Status            string
	AllowGuests       bool
	AllowSpectators   bool
	RoundTimerSeconds int
	CreatedAt         time.Time
	LastActivityAt    time.Time
	FinishedAt        *time.Time
}
