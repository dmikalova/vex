package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Unlocked Gateway
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Omega.
//	Play: Destroy each creature.
func TestUnlockedGateway(t *testing.T) {
	t.Run("destroys every creature", func(t *testing.T) {
		var ally, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(UnlockedGateway),
				// Omega ends the turn, so a draw step follows; give P1 a stocked
				// deck so the refill draws from the deck rather than recycling the
				// discard (which would pull the just-destroyed ally back to hand).
				Deck: ct.Cards(
					ct.Creature(), ct.Creature(), ct.Creature(),
					ct.Creature(), ct.Creature(), ct.Creature(),
				),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.Power(4)))),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4))))},
		})

		h.P1.Play(UnlockedGateway)

		h.Expect(ally).At(ct.Discard)
		h.Expect(foe).At(ct.Discard)
	})
}
