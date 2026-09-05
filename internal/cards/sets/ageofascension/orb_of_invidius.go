package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Orb of Invidius
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	After a creature reaps, stun it.
var OrbOfInvidius = card.New(
	"Orb of Invidius",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 96),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Stun{Target: card.Target.Triggering}),
)
