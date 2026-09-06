package engine

import (
	"strings"
	"testing"
)

// TestForgeKeyNumberBarred covers the Key Imps' constant rule: a card in play
// that bars a key ordinal stops every player from forging that key, whoever
// controls the card, and renders its line.
func TestForgeKeyNumberBarred(t *testing.T) {
	g := NewGame("A", "B", 1)
	if g.forgeKeyNumberBarred(0) {
		t.Error("no barring card in play should leave forging open")
	}

	def := NewCard("Bronze Key Imp", Dis, Creature, Common,
		WithPower(2), WithRestrictions(Restrictions{NoForgeKeyNumber: 1}))
	g.AddToBattleline(def, 1) // opponent controls the Imp; it still bars player 0

	g.State.Aember[0] = 3 * KeyCost
	keysBefore := g.State.Keys[0]
	g.forgeKey(0)
	if g.State.Keys[0] != keysBefore {
		t.Error("the first key should be barred while a NoForgeKeyNumber:1 card is in play")
	}
	if g.State.Aember[0] != 3*KeyCost {
		t.Error("a barred forge should not spend Æmber")
	}

	g.forgeKeyFree(0)
	if g.State.Keys[0] != keysBefore {
		t.Error("a free forge of a barred ordinal should also be barred")
	}

	// A second key is not barred by the first-key Imp.
	g.State.Keys[0] = 1
	g.forgeKey(0)
	if g.State.Keys[0] != 2 {
		t.Error("the second key should forge when only the first is barred")
	}

	if !strings.Contains(RenderCardRules(&def), "Players cannot forge their first key.") {
		t.Error("card rules should render the no-forge line")
	}
}

// TestOrdinalWord covers the ordinal words the Key Imps print.
func TestOrdinalWord(t *testing.T) {
	for n, want := range map[int]string{1: "first", 2: "second", 3: "third", 4: "4th"} {
		if got := ordinalWord(n); got != want {
			t.Errorf("ordinalWord(%d) = %q, want %q", n, got, want)
		}
	}
}
