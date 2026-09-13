package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Quixxle Stone
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	If a player has more Creatures in play than their opponent, they cannot play Creatures.
var QuixxleStone = card.New(
	"Quixxle Stone",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "338"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithCannotPlayWhile(card.ConditionalPlayBar{
		Type: card.Type.Creature,
		When: card.ControlsMoreCreatures{},
	}),
)
