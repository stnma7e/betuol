// Package character implements the interaction and attributes for GameObjects to behave like characters in a game.
// This includes NPC's and players, and other interactable entities.
package character

const (
	HEALTH       = iota
	MANA         = iota
	STRENGTH     = iota
	INTELLIGENCE = iota
	RANGEOFSIGHT = iota

	NUM_ATTRIBUTES = iota

	RESIZESTEP = 20
)

// CharacterAttributes is a helper object to group attributes of characters.
type CharacterAttributes struct {
	Attributes [NUM_ATTRIBUTES]float32
}
