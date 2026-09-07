//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// XenosBloodshadow
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Human • Witch
//
//	Elusive. Hazardous 6. Poison. Skirmish.
var XenosBloodshadow = card.New(
	"Xenos Bloodshadow",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 404),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
