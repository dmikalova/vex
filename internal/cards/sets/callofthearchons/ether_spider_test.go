package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Ether Spider
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Beast
//
//	Ether Spider deals no damage when fighting.
//	If Æmber would be added to your opponent's pool, instead Ether Spider captures it.
func TestEtherSpider(t *testing.T) {
	t.Run("deals no damage when fighting", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(EtherSpider),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
			)},
		})

		h.P1.Fight(EtherSpider, foe)

		h.Expect(foe).Damage(0)
		h.Expect(EtherSpider).Damage(3)
	})

	t.Run("captures Æmber that would be added to the opponent's pool", func(t *testing.T) {
		var spider ct.Card
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(FertilityChant),
				InPlay: ct.Cards(
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&spider, EtherSpider))},
		})

		h.P1.Play(FertilityChant)
		h.P1.Reap(reaper)

		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(2)
		h.Expect(spider).AmberOn(5)
	})
}
