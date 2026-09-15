package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Walls' Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Walls' Blaster to Chief Engineer Walls -> for each upgrade on Chief Engineer Walls, stun a creature."
func TestWallsBlaster(t *testing.T) {
	t.Run("the deal 2 damage branch damages a chosen creature", func(t *testing.T) {
		var carrier, victim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&carrier,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						WallsBlaster,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&victim, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Reap(carrier)
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("deal")
		h.P1.ClickCard(victim)

		h.Expect(victim).Damage(2)
	})

	t.Run(
		"the attach branch homes onto Walls and stuns one creature per upgrade",
		func(t *testing.T) {
			var carrier, walls, foe1, foe2 ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.StarAlliance,
					InPlay: ct.Cards(
						ct.Upgraded(
							ct.Bind(
								&carrier,
								ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
							),
							WallsBlaster,
						),
						ct.Upgraded(
							ct.Bind(&walls, ChiefEngineerWalls),
							ct.Upgrade(ct.OfHouse(card.House.StarAlliance)),
						),
					),
				},
				P2: ct.Side{
					InPlay: ct.Cards(
						ct.Bind(&foe1, ct.Creature(ct.Power(6))),
						ct.Bind(&foe2, ct.Creature(ct.Power(6))),
					),
				},
			})

			h.P1.Reap(carrier)
			h.P1.ExpectPrompt("Choose one")
			h.P1.ClickOption("attach")

			// Walls' Blaster is now the second upgrade on Walls, so two creatures stun.
			h.Expect(WallsBlaster).At(ct.Attached)
			h.P1.ClickCard(foe1)
			h.P1.ClickCard(foe2)

			h.Expect(foe1).Stunned(true)
			h.Expect(foe2).Stunned(true)
		},
	)
}
