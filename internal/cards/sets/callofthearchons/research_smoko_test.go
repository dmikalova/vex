package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Research Smoko
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Mutant
//
//	Destroyed: Archive the top card of your deck.
func TestResearchSmoko(t *testing.T) {
	t.Run("archives the top card of your deck when destroyed", func(t *testing.T) {
		var smoko, top, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&smoko, ResearchSmoko)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				),
			},
		})

		h.P1.Fight(smoko, foe)

		h.Expect(smoko).At(ct.Discard)
		h.Expect(top).At(ct.Archives)
	})
}
