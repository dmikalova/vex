package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Word of Returning
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Deal 1 damage to each enemy creature for each Æmber on it, and move all Æmber from each enemy creature to your pool.
func TestWordOfReturning(t *testing.T) {
	t.Run("damages each enemy creature per Æmber on it, then takes that Æmber", func(t *testing.T) {
		var laden, bare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(WordOfReturning),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&laden, ct.Creature(ct.Power(5))),
					ct.Bind(&bare, ct.Creature(ct.Power(5))),
				),
			},
		})
		h.Game().AddAmberOn(laden.ID(), 2)

		h.P1.Play(WordOfReturning)

		h.Expect(laden).Damage(2).AmberOn(0)
		h.Expect(bare).Damage(0)
		// 1 from the card's own Æmber bonus, 2 returned from the creature.
		h.P1.ExpectAmber(3)
	})
}
