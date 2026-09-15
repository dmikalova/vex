//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Shadowsaurus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Thief
//
//	Action: Move each A from an enemy creature to your opponent's pool. If there was at least 1A on that creature, take control of it. While under your control, it belongs to house Shadows.
var Shadowsaurus = set.New(
	"Shadowsaurus",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "292"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
