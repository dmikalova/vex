//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ScrivenerFavian
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Mutant
//
//	Enhance Capture Capture.
//	When you resolve a PT bonus icon, you may choose to steal 1A instead.
var ScrivenerFavian = set.New(
	"Scrivener Favian",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "155"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
