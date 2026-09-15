package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Peace Accord
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	Play: Each player gains 2 Æmber.
//	After a creature is used to fight, its controller loses 4 Æmber. Destroy Peace Accord.
func TestPeaceAccord(t *testing.T) {
	t.Run("each player gains 2 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Amber: 5,
				Hand:  ct.Cards(PeaceAccord),
			},
			P2: ct.Side{Amber: 0},
		})

		h.P1.Play(PeaceAccord)

		h.P1.ExpectAmber(7)
		h.P2.ExpectAmber(2)
	})

	t.Run(
		"after a creature fights, its controller loses 4 and it is destroyed",
		func(t *testing.T) {
			var attacker, defender ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Brobnar,
					Amber: 5,
					InPlay: ct.Cards(
						PeaceAccord,
						ct.Bind(&attacker, ct.Creature(ct.Power(5))),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(ct.Bind(&defender, ct.Creature(ct.Power(3)))),
				},
			})

			h.P1.Fight(attacker, defender)

			h.P1.ExpectAmber(1)
			h.Expect(PeaceAccord).At(ct.Discard)
		},
	)
}
