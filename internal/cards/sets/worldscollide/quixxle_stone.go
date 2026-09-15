package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Quixxle Stone
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	If a player has more creatures in play than their opponent, they cannot play creatures.
var QuixxleStone = set.New(
	"Quixxle Stone",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "338"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithCannotPlayWhile(card.ConditionalPlayBar{
		Type: card.Type.Creature,
		When: card.ControlsMoreCreatures{},
	}),
)
