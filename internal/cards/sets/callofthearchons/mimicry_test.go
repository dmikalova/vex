package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mimicry
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Play a tactic from your opponent's discard pile.
func TestMimicry(t *testing.T) {
	t.Run("plays an action from the opponent's discard and puts it on top", func(t *testing.T) {
		var copied, buried, prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Mimicry),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&prey, ct.Creature(ct.OfHouse(card.House.Brobnar)))),
				Discard: ct.Cards(
					ct.Bind(&copied, BloodMoney),
					ct.Bind(&buried, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				),
			},
		})

		h.P1.Play(Mimicry)

		// Blood Money's Play exalts an enemy creature twice, resolving under P1's
		// control, and the copied card lands on top of its owner's discard pile.
		h.Expect(prey).At(ct.PlayArea).AmberOn(2)
		h.Expect(copied).At(ct.Discard)
		discard := h.Game().Discard(1)
		if top := discard[len(discard)-1]; top != copied.ID() {
			t.Errorf("top of P2's discard = %d, want the copied card %d", top, copied.ID())
		}
	})

	t.Run("stays in the discard when a play limit blocks it", func(t *testing.T) {
		var copied, filler ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand: ct.Cards(
					ct.Bind(&filler, ct.Creature(ct.OfHouse(card.House.Untamed))),
					Mimicry,
				),
			},
			P2: ct.Side{
				InPlay:  ct.Cards(EmberImp),
				Discard: ct.Cards(ct.Bind(&copied, BloodMoney)),
			},
		})

		// Ember Imp caps P1 at two plays. The filler is play one and Mimicry play
		// two, so the card Mimicry would copy is a third play and is barred; it is
		// left untouched in its owner's discard.
		h.P1.Play(filler)
		h.P1.Play(Mimicry)

		h.Expect(copied).At(ct.Discard)
		if discard := h.Game().Discard(1); len(discard) != 1 || discard[0] != copied.ID() {
			t.Errorf("P2's discard = %v, want just the un-played copy %d", discard, copied.ID())
		}
	})
}
