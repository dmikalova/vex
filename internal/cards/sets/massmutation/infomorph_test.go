package massmutation

import (
	"strings"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// Infomorph
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	Enhance Draw Draw.
func TestInfomorph(t *testing.T) {
	if got := engine.RenderCardText(&Infomorph); !strings.Contains(got, "Enhance Draw Draw.") {
		t.Errorf("Infomorph text = %q, want the Enhance line", got)
	}
}
