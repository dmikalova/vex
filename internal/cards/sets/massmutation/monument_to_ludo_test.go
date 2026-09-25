package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Monument to Ludo
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Action: Move 1 Æmber from a creature to the common supply. If Praefectus Ludo is in your discard pile, move 1 Æmber from the chosen creature to the common supply.
func TestMonumentToLudo(t *testing.T) {
	t.Run("moves 1 Æmber from the chosen creature to the common supply", func(t *testing.T) {
		var rich ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				InPlay: ct.Cards(
					MonumentToLudo,
					ct.Bind(&rich, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature(ct.Power(4)))},
		})
		h.Game().AddAmberOn(rich.ID(), 2)

		h.P1.UseAction(MonumentToLudo)
		h.P1.ClickCard(rich)

		h.Expect(rich).AmberOn(1)
	})

	t.Run(
		"moves 2 Æmber when Praefectus Ludo is in your discard pile",
		func(t *testing.T) {
			var rich ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Saurian,
					InPlay: ct.Cards(
						MonumentToLudo,
						ct.Bind(&rich, ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4))),
					),
					Discard: ct.Cards(ct.Creature(ct.Named(praefectusLudoName))),
				},
				P2: ct.Side{InPlay: ct.Cards(ct.Creature(ct.Power(4)))},
			})
			h.Game().AddAmberOn(rich.ID(), 3)

			h.P1.UseAction(MonumentToLudo)
			h.P1.ClickCard(rich)

			h.Expect(rich).AmberOn(1)
		},
	)
}
