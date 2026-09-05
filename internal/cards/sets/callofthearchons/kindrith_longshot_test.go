package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Kindrith Longshot
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Ranger
//
//	Elusive, Skirmish.
//	Reap: Deal 2 damage to a creature.
func TestKindrithLongshot(t *testing.T) {
	t.Run("deals 2 damage to a chosen creature when it reaps", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(KindrithLongshot)},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Reap(KindrithLongshot)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.PlayArea).Damage(2)
	})
}
