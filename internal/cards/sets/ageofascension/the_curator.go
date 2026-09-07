package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// The Curator
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Scientist
//
//	Friendly artifacts enter play ready.
var TheCurator = card.New(
	"The Curator",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "157"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithFriendlyEntersPlayReady(card.Type.Artifact),
)
