//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CreedOfNature
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Omni: Sacrifice Creed of Nature. If you do, choose a creature. For the remainder of the turn, that creature gains skirmish and assault X. X is its power.
var CreedOfNature = card.New(
	"Creed of Nature",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 385),
	card.WithTraits(card.Traits.Power),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
