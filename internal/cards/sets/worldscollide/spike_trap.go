//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SpikeTrap
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Weapon
//
//	Omni: Sacrifice Spike Trap. If you do, deal 3D to each flank creature.
var SpikeTrap = card.New(
	"Spike Trap",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 261),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Weapon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
