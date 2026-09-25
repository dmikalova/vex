package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Saurus Rex
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Dinosaur • Leader
//
//	Fight/Reap: If Saurus Rex is in the center of your battleline, you may exalt Saurus Rex -> search your deck for a Saurian card, reveal it, and put it into your hand. Shuffle your deck.
func TestSaurusRex(t *testing.T) {
	t.Run("centered: exalts and tutors a Saurian card", func(t *testing.T) {
		var rex, ally, outsider ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&rex, SaurusRex)),
				Deck: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(2))),
					ct.Bind(&outsider, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
				),
			},
		})

		h.P1.Reap(rex)      // a lone creature sits in the center
		h.P1.ClickCard(rex) // accept the optional exalt

		h.Expect(rex).AmberOn(1)
		h.Expect(ally).At(ct.Hand)
		h.Expect(outsider).At(ct.Deck)
	})

	t.Run("off-center: no exalt and no search", func(t *testing.T) {
		var rex, flank, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					ct.Bind(&rex, SaurusRex),
					ct.Bind(&flank, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(2))),
				),
				Deck: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(2))),
				),
			},
		})

		h.P1.Reap(rex) // an even battleline has no center

		h.Expect(rex).AmberOn(0)
		h.Expect(ally).At(ct.Deck)
	})
}
