package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Vezyma Thinkdrone
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Scientist
//
//	Reap: You may archive a friendly creature or artifact from play.
func TestVezymaThinkdrone(t *testing.T) {
	t.Run("may archive a friendly creature from play on reap", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Mars,
				InPlay: ct.Cards(
					VezymaThinkdrone,
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(4))),
				),
			},
		})

		h.P1.Reap(VezymaThinkdrone)
		h.P1.ClickCard(ally)

		h.Expect(ally).At(ct.Archives)
	})
}
