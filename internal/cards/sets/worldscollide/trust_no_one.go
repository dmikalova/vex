package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Trust No One
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: If there are no friendly creatures in play, for each house represented among enemy creatures (to a maximum of 3), steal 1 Æmber. Otherwise, steal 1 Æmber.
var TrustNoOne = card.New(
	"Trust No One",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 248),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				None:   true,
			},
			Then: card.StealAember{
				Amount: 1,
				Per: card.HousesAmong{
					Player: card.Opponent,
					Type:   card.Type.Creature,
					Max:    3,
				},
			},
			Else: card.StealAember{Amount: 1},
		}),
)
