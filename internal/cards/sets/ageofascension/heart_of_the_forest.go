package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Heart of the Forest
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	Each player cannot forge keys while they have more forged keys than their opponent.
var HeartOfTheForest = set.New(
	"Heart of the Forest",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "355"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithRestrictions(card.Restrictions{NoForgeWhileAheadOnKeys: true}),
)
