package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Cyber-Clone
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Play: Purge another creature. Cyber-Clone has power equal to the same creature's printed power and gains its printed armor, keywords, and traits.
var CyberClone = set.New(
	"Cyber-Clone",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "102"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.PurgeCreature{Target: card.Target.Creature.Other()},
				card.CopyPrintedStats{
					Target: card.Target.This,
					Source: card.Target.TheSameCreature,
				},
			},
		}),
)
