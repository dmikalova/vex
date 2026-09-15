//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MutagenicSerum
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Omni: Destroy Mutagenic Serum. You may use friendly Mutant creatures this turn.
var MutagenicSerum = set.New(
	"Mutagenic Serum",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "091"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
