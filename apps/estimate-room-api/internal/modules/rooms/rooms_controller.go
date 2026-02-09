package rooms

import (
	"net/http"
)

type RoomsController any

type roomsController struct{}

func NewRoomsController() RoomsController {
	return &roomsController{}
}

func (c *roomsController) ListRooms(w http.ResponseWriter, r *http.Request) {}
