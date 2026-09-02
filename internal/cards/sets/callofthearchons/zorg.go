package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Zorg
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Beast
//
//	Zorg enters play stunned.
//	Before Fight: Stun the creature Zorg fights and each of its neighbors.
var Zorg = card.New(
	"Zorg",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 206),
	card.WithPower(7),
	card.WithTraits("Beast"),
	card.WithEntersPlay(card.Stun{Target: card.Target.This}),
	card.WithAbility(
		card.Trigger.BeforeFight,
		card.Stun{Target: card.Target.CreatureFought.AndNeighbors()}),
)
