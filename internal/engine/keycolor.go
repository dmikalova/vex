package engine

// KeyColor is the colour of a forged key. KeyForge keys are all alike, but Vactrol
// lets a player pick a colour as they forge, so each of a player's up-to-three
// keys is a distinct colour. KeyColorNone is the zero value (an unforged slot).
type KeyColor uint8

// The key colours a player can forge; KeyColorNone is the unforged zero value.
const (
	KeyColorNone KeyColor = iota
	KeyColorRed
	KeyColorBlue
	KeyColorYellow
)

// keyColorNames maps a KeyColor to its printed name, indexed by the enum value.
var keyColorNames = [...]string{"None", "Red", "Blue", "Yellow"}

// String returns the printed colour name.
func (c KeyColor) String() string {
	if int(c) < len(keyColorNames) {
		return keyColorNames[c]
	}
	return "Unknown"
}

// keyColorOrder is the canonical order colours are offered and displayed in.
var keyColorOrder = [...]KeyColor{KeyColorRed, KeyColorBlue, KeyColorYellow}

// KeyColorPrompt is the option-chooser prompt shown when forging asks which key
// colour to forge. It is exported so a frontend (or test harness) can recognise
// this specific prompt.
const KeyColorPrompt = "Choose a key colour to forge"
