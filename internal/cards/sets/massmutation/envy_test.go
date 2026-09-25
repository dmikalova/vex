package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Envy
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Demon • Sin
//
//	Elusive.
//	Reap: If there are 2 or more friendly Sin creatures in play, Envy captures all your opponent's Æmber.
func TestEnvy(t *testing.T) {
	t.Run(
		"captures all opponent Æmber with two or more friendly Sin creatures",
		func(t *testing.T) {
			var envy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					InPlay: ct.Cards(
						ct.Bind(&envy, Envy),
						ct.Creature(
							ct.OfHouse(card.House.Dis),
							ct.Traits(card.Traits.Sin),
							ct.Power(3),
						),
					),
				},
				P2: ct.Side{Amber: 4},
			})

			h.P1.Reap(envy)

			h.P2.ExpectAmber(0)
			h.Expect(envy).AmberOn(4)
		},
	)

	t.Run("captures nothing with only one friendly Sin creature", func(t *testing.T) {
		var envy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&envy, Envy)),
			},
			P2: ct.Side{Amber: 4},
		})

		h.P1.Reap(envy)

		h.P2.ExpectAmber(4)
		h.Expect(envy).AmberOn(0)
	})
}
