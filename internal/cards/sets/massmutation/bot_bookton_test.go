package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Bot Bookton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant • Scientist
//
//	Reap: Play the top card of your deck.
func TestBotBookton(t *testing.T) {
	t.Run("plays the top card of your deck when it reaps", func(t *testing.T) {
		var top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(BotBookton),
				Deck: ct.Cards(
					ct.Bind(&top, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(BotBookton)

		h.Expect(top).At(ct.PlayArea)
	})
}
