package roomsmodels

type DeckID string

const (
	DeckIDFibonacci DeckID = "FIBONACCI"
	DeckIDTShirt    DeckID = "TSHIRT"
	DeckIDCustom    DeckID = "CUSTOM"
)

func (d DeckID) IsValid() bool {
	switch d {
	case DeckIDFibonacci, DeckIDTShirt, DeckIDCustom:
		return true
	default:
		return false
	}
}
