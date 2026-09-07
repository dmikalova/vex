package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Kangaphant
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Beast
//
//	Each creature gains, "Reap: Destroy this creature."
var Kangaphant = card.New(
	"Kangaphant",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "376"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect:  card.Destroy{Target: card.Target.This},
		}},
	}),
)
