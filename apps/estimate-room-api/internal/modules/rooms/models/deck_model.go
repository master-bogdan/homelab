package roomsmodels

import "strings"

type RoomDeck struct {
	Name   string   `json:"name"`
	Kind   string   `json:"kind"`
	Values []string `json:"values"`
}

func DefaultRoomDeck() RoomDeck {
	return RoomDeck{
		Name: "Fibonacci",
		Kind: "FIBONACCI",
		Values: []string{
			"0",
			"1",
			"2",
			"3",
			"5",
			"8",
			"13",
			"21",
			"?",
		},
	}
}

func (d RoomDeck) IsZero() bool {
	return strings.TrimSpace(d.Name) == "" &&
		strings.TrimSpace(d.Kind) == "" &&
		len(d.Values) == 0
}

func (d RoomDeck) IsValid() bool {
	if strings.TrimSpace(d.Name) == "" {
		return false
	}
	if strings.TrimSpace(d.Kind) == "" {
		return false
	}
	if len(d.Values) == 0 {
		return false
	}

	for _, value := range d.Values {
		if strings.TrimSpace(value) == "" {
			return false
		}
	}

	return true
}
