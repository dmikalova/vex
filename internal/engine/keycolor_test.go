package engine

import (
	"reflect"
	"testing"
)

func TestKeyColorString(t *testing.T) {
	cases := map[KeyColor]string{
		KeyColorNone: "None", KeyColorRed: "Red", KeyColorBlue: "Blue",
		KeyColorYellow: "Yellow", KeyColor(9): "Unknown",
	}
	for c, want := range cases {
		if got := c.String(); got != want {
			t.Errorf("KeyColor(%d).String() = %q, want %q", c, got, want)
		}
	}
}

// TestForgeRecordsKeyColor forges all three keys and checks the chosen colours are
// recorded in forge order: the default chooser picks the first remaining colour
// each time, and the last key's colour is forced.
func TestForgeRecordsKeyColor(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 3 * KeyCost
	for i := 0; i < 3; i++ {
		g.forgeKey(0)
	}
	got := g.KeyColors(0)
	want := []KeyColor{KeyColorRed, KeyColorBlue, KeyColorYellow}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("KeyColors = %v, want %v", got, want)
	}
}

// TestForgeKeyColorChoice lets the player pick a non-default colour: choosing
// index 1 each prompt yields Blue then Yellow, then Red is forced.
func TestForgeKeyColorChoice(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetChooser(0, optionPicker{idx: 1})
	g.State.Aember[0] = 3 * KeyCost
	for i := 0; i < 3; i++ {
		g.forgeKey(0)
	}
	got := g.KeyColors(0)
	want := []KeyColor{KeyColorBlue, KeyColorYellow, KeyColorRed}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("KeyColors = %v, want %v", got, want)
	}
}

// TestForgeBeyondKeysLeavesColorsUntouched covers forging when no colour remains:
// a fourth forge records nothing rather than overrunning the colour slots.
func TestForgeBeyondKeysLeavesColorsUntouched(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 4 * KeyCost
	for i := 0; i < 4; i++ {
		g.forgeKey(0)
	}
	if got := g.State.KeyColors[0]; got != [KeysToWin]KeyColor{
		KeyColorRed,
		KeyColorBlue,
		KeyColorYellow,
	} {
		t.Errorf("key colours after 4 forges = %v, want [Red Blue Yellow]", got)
	}
}
