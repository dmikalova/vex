package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Trust No One
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose one:
//	- If there are no friendly creatures in play, for each house represented among enemy creatures, steal 1 Æmber
//	- Steal 1 Æmber.
var TrustNoOne = set.New(
	"Trust No One",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "248"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseOne{
			Options: []card.Effect{
				card.Conditional{
					Cond: card.CardsInPlay{
						Player: card.Controller,
						Type:   card.Type.Creature,
						None:   true,
					},
					Then: card.StealAember{
						Amount: 1,
						Per: card.HousesAmong{
							Player: card.Opponent,
							Type:   card.Type.Creature,
						},
					},
				},
				card.StealAember{Amount: 1},
			},
		}),
)
