package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Exterminate! Exterminate!
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy each non-Mars Creature with power less than the number of friendly Mars Creatures you control.
var ExterminateExterminate = card.New(
	"Exterminate! Exterminate!",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "180"),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.
				ExceptHouse(card.House.Self).
				Refine(card.PowerLessThan(card.InPlay{
					Player: card.Controller,
					Type:   card.Type.Creature,
					House:  card.House.Self,
				})),
		}),
)
