package massmutation

import (
	"strings"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// Splinter
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant • Thief
//
//	Enhance Damage Damage Damage Damage Damage Damage.
func TestSplinter(t *testing.T) {
	if got := engine.RenderCardText(&Splinter); !strings.Contains(
		got, "Enhance Damage Damage Damage Damage Damage Damage.",
	) {
		t.Errorf("Splinter text = %q, want the Enhance line", got)
	}
}
