package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Phoenix Heart
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains, "Destroyed: Put this creature into its owner's hand, and deal 3 damage to each creature."
func TestPhoenixHeart(t *testing.T) {
	t.Run(
		"returns its host to hand and deals 3 to each creature when the host is destroyed",
		func(t *testing.T) {
			var host, bystander, enemy ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Brobnar,
					InPlay: ct.Cards(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
						ct.Bind(
							&bystander,
							ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5)),
						),
					),
					Hand: ct.Cards(PhoenixHeart),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(5))),
					),
				},
			})

			h.P1.Play(PhoenixHeart)
			h.P1.ClickCard(host) // choose which creature to upgrade

			h.Game().DestroyEach(0, []engine.LocalID{host.ID()})

			h.Expect(host).At(ct.Hand) // returned to its owner's hand
			h.Expect(bystander).At(ct.PlayArea).Damage(3)
			h.Expect(enemy).Damage(3)
		},
	)
}
