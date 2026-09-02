package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Oath of Poverty
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each friendly artifact. For each card destroyed this way, gain 2 Æmber.
func TestOathOfPoverty(t *testing.T) {
	t.Run("destroys your artifacts and gains 2 Æmber for each", func(t *testing.T) {
		var mine, theirs ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(OathOfPoverty),
				InPlay: ct.Cards(
					ct.Bind(&mine, ct.Artifact(ct.OfHouse(card.House.Sanctum))),
					ct.Artifact(ct.OfHouse(card.House.Sanctum)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&theirs, ct.Artifact(ct.OfHouse(card.House.Mars))),
			)},
		})

		h.P1.Play(OathOfPoverty)

		h.P1.ExpectAmber(5) // 1 Æmber bonus plus 2 for each of the two destroyed
		h.Expect(mine).At(ct.Discard)
		h.Expect(theirs).At(ct.PlayArea)
	})

	t.Run("gains nothing when you control no artifacts", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(OathOfPoverty)},
		})

		h.P1.Play(OathOfPoverty)

		h.P1.ExpectAmber(1)
	})
}
