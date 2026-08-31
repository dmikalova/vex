package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Poison Wave
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Deal 2 damage to each creature.
func TestPoisonWave(t *testing.T) {
	t.Run("deals 2 damage to each creature", func(t *testing.T) {
		var ally, foe, weakFoe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(PoisonWave),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Shadows), ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
					ct.Bind(&weakFoe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Play(PoisonWave)

		h.Expect(ally).At(ct.PlayArea).Damage(2)
		h.Expect(foe).At(ct.PlayArea).Damage(2)
		h.Expect(weakFoe).At(ct.Discard)
	})
}
