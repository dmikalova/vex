//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DarkAemberVault
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Special
//	Traits: Location
//
//	After you play a Mutant creature, draw a card.
//	Each friendly Mutant creature gets +2 power.
var DarkAemberVault = set.New(
	"Dark Aember Vault",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Special,
	card.Provenance(card.MM, "001"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
