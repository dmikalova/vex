package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Deep Probe
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a house. Reveal your opponent's hand. Discard each creature of the chosen house from your opponent's hand.
func TestDeepProbe(t *testing.T) {
	t.Run(
		"discards each creature of the chosen house from the opponent's hand",
		func(t *testing.T) {
			var marsCreature, marsAction, sanctumCreature ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Mars,
					Hand:  ct.Cards(DeepProbe),
				},
				P2: ct.Side{
					Hand: ct.Cards(
						ct.Bind(
							&marsCreature,
							ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)),
						),
						ct.Bind(&marsAction, ct.Tactic(ct.OfHouse(card.House.Mars))),
						ct.Bind(
							&sanctumCreature,
							ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(2)),
						),
					),
				},
			})

			h.P1.Play(DeepProbe)
			h.P1.ExpectPrompt("Choose a house").Source("Deep Probe")
			h.P1.ClickOption("Mars")

			h.Expect(marsCreature).At(ct.Discard) // a Mars creature is discarded
			h.Expect(marsAction).At(ct.Hand)      // a Mars action is not a creature
			h.Expect(sanctumCreature).At(ct.Hand) // a Sanctum creature is the wrong house
		},
	)
}
