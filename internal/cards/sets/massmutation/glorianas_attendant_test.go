package massmutation

import (
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Gloriana's Attendant
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Mutant
//
//	Enhance Æmber Æmber.
func TestGlorianasAttendant(t *testing.T) {
	if got := engine.RenderCardText(&GlorianasAttendant); !strings.Contains(
		got, "Enhance Æmber Æmber.",
	) {
		t.Errorf("Gloriana's Attendant text = %q, want the Enhance line", got)
	}
}
