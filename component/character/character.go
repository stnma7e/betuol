package character

import (
	"smig/component"
)

type CharacterAttributes struct {
	Attributes [NUM_ATTRIBUTES]float32
	Description, Greeting string
}

type Character struct {
	Id component.GOiD
}

func (ca *CharacterAttributes) Greet() string {
	return ca.Description + " says: " + ca.Greeting
}