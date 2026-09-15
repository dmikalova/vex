package massmutation

import (
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Mutant Cutpurse
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Thief
//
//	Enhance Damage Damage Damage.
func TestMutantCutpurse(t *testing.T) {
	if got := engine.RenderCardText(&MutantCutpurse); !strings.Contains(
		got, "Enhance Damage Damage Damage.",
	) {
		t.Errorf("Mutant Cutpurse text = %q, want the Enhance line", got)
	}
}
