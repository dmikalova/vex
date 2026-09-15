package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Soul Snatcher
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Vehicle
//
//	Each creature gains, "Destroyed: Gain 1 Æmber."
func TestSoulSnatcher(t *testing.T) {
	t.Run("a destroyed creature's owner gains 1 Æmber", func(t *testing.T) {
		var attacker, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				InPlay: ct.Cards(
					SoulSnatcher,
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(5))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2))),
				),
			},
		})

		h.P1.Fight(attacker, foe)

		h.Expect(foe).At(ct.Discard)
		h.P2.ExpectAmber(1) // the destroyed creature's owner
	})
}
