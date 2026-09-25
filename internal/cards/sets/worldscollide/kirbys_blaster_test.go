package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Kirby's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Kirby's Blaster to Com. Officer Kirby -> draw 2 cards."
func TestKirbysBlaster(t *testing.T) {
	t.Run("draws 2 cards on the attach payoff", func(t *testing.T) {
		var carrier, top1, top2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&carrier,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						KirbysBlaster,
					),
					ct.Creature(
						ct.OfHouse(card.House.StarAlliance),
						ct.Named("Com. Officer Kirby"),
					),
				),
				Deck: ct.Cards(
					ct.Bind(&top1, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3))),
					ct.Bind(&top2, ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(carrier)
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("attach")

		h.Expect(top1).At(ct.Hand)
		h.Expect(top2).At(ct.Hand)
	})
}
