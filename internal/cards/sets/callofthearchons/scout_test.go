package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Scout
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Give skirmish to, ready, and fight with up to 2 different friendly creatures, one at a time.
func TestScout(t *testing.T) {
	brobnar := ct.OfHouse(card.House.Brobnar)

	t.Run("grants skirmish so a chosen creature takes no fight damage", func(t *testing.T) {
		var scout, ogre ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(Scout),
				InPlay: ct.Cards(ct.Bind(&scout, ct.Creature(brobnar, ct.Power(3)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&ogre, ct.Creature(brobnar, ct.Power(5)))),
			},
		})

		h.P1.Play(Scout)
		h.P1.ClickCard(scout) // the lone enemy is fought automatically

		// Skirmish means the attacker takes no damage back from the 5-power ogre.
		h.Expect(scout).At(ct.PlayArea).Damage(0).Exhausted()
		h.Expect(ogre).At(ct.PlayArea).Damage(3)
	})

	t.Run("readies and fights with up to 2 creatures, one at a time", func(t *testing.T) {
		var wolf, bear, wall ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(Scout),
				InPlay: ct.Cards(
					ct.Bind(&wolf, ct.Creature(brobnar, ct.Power(3))),
					ct.Bind(&bear, ct.Creature(brobnar, ct.Power(4))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&wall, ct.Creature(brobnar, ct.Power(20)))),
			},
		})
		wolf.Exhaust()
		bear.Exhaust()

		h.P1.Play(Scout)
		h.P1.ClickCard(wolf)
		h.P1.ClickCard(bear)

		// Both readied, fought the wall for 3 + 4, and took no damage back.
		h.Expect(wall).At(ct.PlayArea).Damage(7)
		h.Expect(wolf).At(ct.PlayArea).Damage(0).Exhausted()
		h.Expect(bear).At(ct.PlayArea).Damage(0).Exhausted()
	})
}
