package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Ambassador Liu
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant • Politician
//
//	Action: Discard a card from your hand. If it is a Dis or Shadows card, steal 1 Æmber. If it is a Logos or Untamed card, gain 2 Æmber. If it is a Sanctum or Saurian card, Ambassador Liu captures 3 Æmber from your opponent.
func TestAmbassadorLiu(t *testing.T) {
	t.Run("a discarded Shadows card steals 1 Æmber", func(t *testing.T) {
		var pitched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(AmbassadorLiu),
				Hand:   ct.Cards(ct.Bind(&pitched, ct.Creature(ct.OfHouse(card.House.Shadows)))),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(AmbassadorLiu)

		h.Expect(pitched).At(ct.Discard)
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(2)
	})

	t.Run("a discarded Untamed card gains 2 Æmber", func(t *testing.T) {
		var pitched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(AmbassadorLiu),
				Hand:   ct.Cards(ct.Bind(&pitched, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.UseAction(AmbassadorLiu)

		h.Expect(pitched).At(ct.Discard)
		h.P1.ExpectAmber(2)
	})

	t.Run("a discarded Sanctum card captures 3 Æmber onto Ambassador Liu", func(t *testing.T) {
		var liu, pitched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&liu, AmbassadorLiu)),
				Hand:   ct.Cards(ct.Bind(&pitched, ct.Creature(ct.OfHouse(card.House.Sanctum)))),
			},
			P2: ct.Side{Amber: 4},
		})

		h.P1.UseAction(AmbassadorLiu)

		h.Expect(pitched).At(ct.Discard)
		h.Expect(liu).AmberOn(3)
		h.P2.ExpectAmber(1)
	})

	t.Run("a discarded Star Alliance card does nothing", func(t *testing.T) {
		var pitched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(AmbassadorLiu),
				Hand: ct.Cards(
					ct.Bind(&pitched, ct.Creature(ct.OfHouse(card.House.StarAlliance))),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.UseAction(AmbassadorLiu)

		h.Expect(pitched).At(ct.Discard)
		h.P1.ExpectAmber(0)
		h.P2.ExpectAmber(3)
	})
}
