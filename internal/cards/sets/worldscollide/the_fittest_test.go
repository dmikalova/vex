package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Fittest
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Give each friendly creature a +1 power counter.
func TestTheFittest(t *testing.T) {
	t.Run("gives each friendly creature a +1 power counter", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(TheFittest),
				InPlay: ct.Cards(
					ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.Play(TheFittest)

		h.Expect(ally).Power(4)
	})
}
