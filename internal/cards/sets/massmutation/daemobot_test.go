package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Daemo-Bot
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Scientist
//
//	Reap: Discard a card from your hand -> draw a card.
//	Destroyed: Steal 1 Æmber.
func TestDaemoBot(t *testing.T) {
	t.Run("discards a card then draws when it reaps", func(t *testing.T) {
		var handCard, deckCard ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(DaemoBot),
				Hand:   ct.Cards(ct.Bind(&handCard, ct.Creature())),
				Deck:   ct.Cards(ct.Bind(&deckCard, ct.Creature())),
			},
		})

		h.P1.Reap(DaemoBot)

		h.Expect(handCard).At(ct.Discard)
		h.Expect(deckCard).At(ct.Hand)
	})

	t.Run("steals 1 Æmber when destroyed", func(t *testing.T) {
		var bot, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&bot, DaemoBot)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
				),
				Amber: 2,
			},
		})

		h.P1.Fight(bot, foe)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(1)
	})
}
