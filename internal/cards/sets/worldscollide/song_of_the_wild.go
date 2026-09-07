package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Song of the Wild
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Each friendly creature gains, "Reap: Gain 1 Æmber."
var SongOfTheWild = card.New(
	"Song of the Wild",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "364"),
	card.WithAbility(
		card.Trigger.Play, card.GainAbility{
			Target: card.Target.EachFriendlyCreature,
			Ability: card.Ability{
				Trigger: card.Trigger.Reap,
				Effect: card.GainAember{
					Player: card.Controller,
					Amount: 1,
				},
			},
		}),
)
