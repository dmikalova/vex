package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Tribute
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: The most powerful friendly creature captures 2 Æmber from your opponent. You may exalt the chosen creature to repeat the preceding effect.
func TestTribute(t *testing.T) {
	t.Run("the most powerful friendly creature captures 2", func(t *testing.T) {
		var big ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(Tribute),
				InPlay: ct.Cards(
					ct.Bind(&big, ct.Creature(
						ct.OfHouse(card.House.Saurian),
						ct.Power(9),
					)),
					ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4)),
				),
			},
			P2: ct.Side{Amber: 6},
		})

		h.P1.Play(Tribute)
		h.P1.ClickOption("No") // decline the exalt

		h.Expect(big).AmberOn(2)
		h.P2.ExpectAmber(4)
	})

	t.Run("exalting the chosen creature repeats the capture", func(t *testing.T) {
		var big ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(Tribute),
				InPlay: ct.Cards(
					ct.Bind(&big, ct.Creature(
						ct.OfHouse(card.House.Saurian),
						ct.Power(9),
					)),
					ct.Creature(ct.OfHouse(card.House.Saurian), ct.Power(4)),
				),
			},
			P2: ct.Side{Amber: 6},
		})

		h.P1.Play(Tribute)
		h.P1.ClickOption("Yes") // exalt the chosen creature to repeat

		h.Expect(big).AmberOn(5) // 2 captured + 1 exalted + 2 captured
		h.P2.ExpectAmber(2)
	})

	t.Run(
		"a different creature may capture on the repeat, exalt stays on the first",
		func(t *testing.T) {
			var first, second ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Saurian,
					Hand:  ct.Cards(Tribute),
					InPlay: ct.Cards(
						ct.Bind(&first, ct.Creature(
							ct.OfHouse(card.House.Saurian),
							ct.Power(9),
						)),
						ct.Bind(&second, ct.Creature(
							ct.OfHouse(card.House.Saurian),
							ct.Power(9),
						)),
					),
				},
				P2: ct.Side{Amber: 6},
			})

			h.P1.Play(Tribute)
			h.P1.ClickCard(first)   // tie-break: first captures 2
			h.P1.ClickOption("Yes") // exalt that same creature to repeat
			h.P1.ClickCard(second)  // tie-break on the repeat: a different creature captures

			h.Expect(first).AmberOn(3)  // 2 captured + 1 exalted
			h.Expect(second).AmberOn(2) // 2 captured on the repeat
			h.P2.ExpectAmber(2)
		},
	)
}
