package engine

import (
	"reflect"
	"testing"
)

func TestKeyColorString(t *testing.T) {
	cases := map[KeyColor]string{
		KeyColorNone: "None", KeyColorRed: "Red", KeyColorBlue: "Blue",
		KeyColorYellow: "Yellow", KeyColorColorless: "Colorless", KeyColor(9): "Unknown",
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

// TestForgeFourthKeyIsColorless covers forging past the three palette colours: the
// fourth key has no colour left to pick, so it is recorded as KeyColorColorless
// rather than overrunning or reusing a colour.
func TestForgeFourthKeyIsColorless(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.State.Aember[0] = 4 * KeyCost
	for i := 0; i < 4; i++ {
		g.forgeKey(0)
	}
	if got := g.Keys(0); got != MaxKeys {
		t.Errorf("keys after 4 forges = %d, want %d", got, MaxKeys)
	}
	if got := g.State.KeyColors[0]; got != [MaxKeys]KeyColor{
		KeyColorRed,
		KeyColorBlue,
		KeyColorYellow,
		KeyColorColorless,
	} {
		t.Errorf("key colours after 4 forges = %v, want [Red Blue Yellow Colorless]", got)
	}
}
