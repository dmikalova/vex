package massmutation

import (
	"strings"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// Mutagenesis Researcher
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant • Scientist
//
//	Enhance Æmber Capture Damage Draw.
func TestMutagenesisResearcher(t *testing.T) {
	if got := engine.RenderCardText(&MutagenesisResearcher); !strings.Contains(
		got, "Enhance Æmber Capture Damage Draw.",
	) {
		t.Errorf("Mutagenesis Researcher text = %q, want the Enhance line", got)
	}
}
