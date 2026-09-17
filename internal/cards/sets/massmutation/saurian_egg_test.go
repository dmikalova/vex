package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Saurian Egg
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Armor:  5
//	Bonus:  Æmber
//	Traits: Dinosaur • Egg
//
//	Versatile.
//	Saurian Egg cannot fight.
//	Saurian Egg cannot reap.
//	Action: Discard the top 2 cards of your deck. For each Saurian creature discarded this way, put it into play ready. Give it three +1 power counters. If you discard a Saurian creature this way, destroy Saurian Egg.
func TestSaurianEgg(t *testing.T) {
	t.Run(
		"hatches a Saurian creature discarded this way, then destroys itself",
		func(t *testing.T) {
			var egg, saur, filler ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Saurian,
					InPlay: ct.Cards(ct.Bind(&egg, SaurianEgg)),
					Deck: ct.Cards(
						ct.Bind(&saur, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(2))),
						ct.Bind(&filler, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(4))),
					),
				},
			})

			h.P1.UseAction(SaurianEgg)

			// The Saurian creature hatches ready with three +1 power counters (2 -> 5).
			h.Expect(saur).At(ct.PlayArea).Ready().Power(5)
			// A non-Saurian card discarded this way is not reanimated.
			h.Expect(filler).At(ct.Discard)
			// Hatching a Saurian creature destroys the egg.
			h.Expect(egg).At(ct.Discard)
		},
	)

	t.Run("with no Saurian creature discarded the egg survives", func(t *testing.T) {
		var egg, top, next ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&egg, SaurianEgg)),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(4))),
					ct.Bind(&next, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				),
			},
		})

		h.P1.UseAction(SaurianEgg)

		h.Expect(top).At(ct.Discard)
		h.Expect(next).At(ct.Discard)
		h.Expect(egg).At(ct.PlayArea)
	})
}
