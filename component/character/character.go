package character

import (
	"betuol/component"
)

const (
	HEALTH       = iota
	MANA         = iota
	STRENGTH     = iota
	INTELLIGENCE = iota
	RANGEOFSIGHT = iota

	NUM_ATTRIBUTES = iota

	RESIZESTEP = 20
)

type Interaction func(id1 component.GOiD, id2 component.GOiD)

type CharacterAttributes struct {
	Attributes            [NUM_ATTRIBUTES]float32
	Description, Greeting string
}

func (ca *CharacterAttributes) Greet() string {
	return ca.Description + " says: " + ca.Greeting
}
