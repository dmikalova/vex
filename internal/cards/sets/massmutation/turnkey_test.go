package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Turnkey
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Demon
//
//	Play: Unforge one of your opponent's keys -> when Turnkey leaves play, your opponent forges a key at no cost.
func TestTurnkey(t *testing.T) {
	t.Run("unforges a key, then hands it back when it leaves play", func(t *testing.T) {
		var turnkey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&turnkey, Turnkey)),
			},
			P2: ct.Side{
				ForgedKeys: []engine.KeyColor{card.KeyColor.Red},
			},
		})

		h.P1.Play(turnkey)
		h.P2.ExpectKeys(0)

		h.Game().DestroyEach(0, []engine.LocalID{turnkey.ID()})
		h.P2.ExpectKeys(1)
	})

	t.Run("does not arm the forge when the opponent had no key to unforge", func(t *testing.T) {
		var turnkey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(ct.Bind(&turnkey, Turnkey)),
			},
		})

		h.P1.Play(turnkey)
		h.P2.ExpectKeys(0)

		h.Game().DestroyEach(0, []engine.LocalID{turnkey.ID()})
		h.P2.ExpectKeys(0)
	})
}
