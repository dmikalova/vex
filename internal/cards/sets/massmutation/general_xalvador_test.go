package massmutation

import (
	"strings"
	"testing"

	"github.com/dmikalova/vex/internal/engine"
)

// General Xalvador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Enhance Capture Capture.
func TestGeneralXalvador(t *testing.T) {
	if got := engine.RenderCardText(&GeneralXalvador); !strings.Contains(
		got, "Enhance Capture Capture.",
	) {
		t.Errorf("General Xalvador text = %q, want the Enhance line", got)
	}
}
