package cards

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
)

// TestEveryCardDeclaresItsSet guards the rule that deck generation never infers a
// card's home set from provenance (ADR 0003): every registered card must declare
// its set explicitly, which a set package's registrar (set.New) stamps for it. A
// card built through the bare card.New without InSet registers with the zero set
// and fails here — the failure it would otherwise cause (a card silently absent
// from its set's pool) is caught at the source instead.
func TestEveryCardDeclaresItsSet(t *testing.T) {
	for i := range card.Cards() {
		rc := card.Cards()[i]
		if rc.Set.Name == "" {
			t.Errorf(
				"card %q registered with no set: author it through its set's set.New",
				rc.Def.Name,
			)
		}
	}
}
