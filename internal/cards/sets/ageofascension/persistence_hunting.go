package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Persistence Hunting
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose a house. Exhaust each enemy creature of the chosen house.
var PersistenceHunting = set.New(
	"Persistence Hunting",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "328"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Exhaust{Target: card.Target.EachEnemyCreature.House(card.Houses.Chosen)},
		}),
)
