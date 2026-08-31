package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Feeding Pit
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Discard a creature from your hand -> gain 1 Æmber.
func TestFeedingPit(t *testing.T) {
	t.Run("discards a creature from hand to gain 1 Æmber", func(t *testing.T) {
		var beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(FeedingPit),
				Hand:   ct.Cards(ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Mars)))),
			},
		})

		h.P1.UseAction(FeedingPit)

		h.Expect(beast).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})

	t.Run("gains nothing with no creature to discard", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(FeedingPit),
				Hand:   ct.Cards(ct.Tactic(ct.OfHouse(card.House.Mars))),
			},
		})

		h.P1.UseAction(FeedingPit)

		h.P1.ExpectAmber(0)
	})
}
