package callofthearchons

import (
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// Protect the Weak
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature gains +1 armor and taunt.
func TestProtectTheWeak(t *testing.T) {
	if got := engine.RenderCardText(
		&ProtectTheWeak,
	); !strings.Contains(
		got,
		"This creature gains +1 armor and taunt.",
	) {
		t.Errorf("Protect the Weak text = %q, want +1 armor and taunt", got)
	}
}
