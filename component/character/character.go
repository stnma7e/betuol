package character

type CharacterAttributes struct {
	Attributes [NUM_ATTRIBUTES]float32
	Description, Greeting string
}

func (ca *CharacterAttributes) Greet() string {
	return ca.Description + " says: " + ca.Greeting
}