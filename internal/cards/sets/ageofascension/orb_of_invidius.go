package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Orb of Invidius
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	After a creature reaps, stun it.
var OrbOfInvidius = set.New(
	"Orb of Invidius",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "96"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Stun{Target: card.Target.Triggering}),
)
