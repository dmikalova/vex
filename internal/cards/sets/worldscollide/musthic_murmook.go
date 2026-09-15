package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Musthic Murmook
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast
//
//	Each player's keys cost +1 Æmber.
//	Play: Deal 4 damage to a creature.
var MusthicMurmook = set.New(
	"Musthic Murmook",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "361"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Beast),
	card.WithKeyCost(card.KeyCostChange(card.EachPlayer, 1)),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 4,
			Target: card.Target.Creature,
		}),
)
