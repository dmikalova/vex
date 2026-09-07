//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Terrordactyl
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  12
//	Traits: Beast
//
//	Terrordactyl enters play stunned.
//	Terrordactyl only deals 4D when fighting.
//	Before Fight: Deal 4D to each neighbor of the creature Terrordactyl fights.
var Terrordactyl = card.New(
	"Terrordactyl",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 211),
	card.WithPower(12),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
