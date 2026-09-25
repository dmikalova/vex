package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Miasma Bomb
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Action: Destroy Miasma Bomb -> your opponent skips the "forge a key" phase during their next turn.
//	Enhance Damage.
func TestMiasmaBomb(t *testing.T) {
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Shadows,
			InPlay: ct.Cards(MiasmaBomb),
		},
	})

	h.P1.UseAction(MiasmaBomb)

	h.Expect(MiasmaBomb).At(ct.Discard)
	if !h.Game().State.SkipForgeNext[1].Value {
		t.Error("the opponent should be set to skip their next forge step")
	}
}
