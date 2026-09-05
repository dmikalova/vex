package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Zyzzix the Many
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Soldier
//
//	Fight/Reap: You may reveal a creature from your hand and archive it -> give Zyzzix the Many 3 +1 power counters.
func TestZyzzixTheMany(t *testing.T) {
	t.Run("reaping archives a revealed creature and grows Zyzzix", func(t *testing.T) {
		var zyzzix, creature ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&zyzzix, ZyzzixTheMany)),
				Hand: ct.Cards(
					ct.Bind(&creature, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(zyzzix)
		h.P1.ClickCard(creature)

		h.Expect(creature).At(ct.Archives)
		h.Expect(zyzzix).Power(6)
	})

	t.Run("with no creature in hand Zyzzix does not grow", func(t *testing.T) {
		var zyzzix, tactic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&zyzzix, ZyzzixTheMany)),
				Hand:   ct.Cards(ct.Bind(&tactic, ct.Tactic(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.Reap(zyzzix)

		h.Expect(tactic).At(ct.Hand)
		h.Expect(zyzzix).Power(3)
	})
}
