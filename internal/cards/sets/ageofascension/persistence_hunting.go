package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Persistence Hunting
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a house - exhaust each enemy creature of the chosen house.
var PersistenceHunting = card.New(
	"Persistence Hunting",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 328),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Exhaust{Target: card.Target.EachEnemyCreature.OfChosenHouse()},
		}),
)
