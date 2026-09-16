package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Citizen Shrix
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant
//
//	Play/Reap: Exalt Citizen Shrix. Steal 1 Æmber.
var CitizenShrix = set.New(
	"Citizen Shrix",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "186"),
	card.InCluster(card.Pulled(monumentToShrixCluster, 1, 1)),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.PlayReap, card.Sentences{Effects: []card.Effect{
			card.Exalt{Target: card.Target.This, Amount: 1},
			card.StealAember{Amount: 1},
		}}),
)
