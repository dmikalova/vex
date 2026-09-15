//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BlossomDrake
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dragon
//
//	Blossom Drake gets +1 power for each artifact in play.
//	Each artifact's text box is considered blank (except for traits).
var BlossomDrake = set.New(
	"Blossom Drake",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "395"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dragon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
