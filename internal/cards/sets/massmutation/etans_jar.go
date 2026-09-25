package massmutation

import "github.com/dmikalova/vex/internal/card"

// Etan's Jar
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Play: Name a card. Cards with that name cannot be played until Etan's Jar leaves play.
var EtansJar = set.New(
	"Etan's Jar",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "035"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Play, card.NameCard{}),
)
