//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CreedOfNurture
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Omni: Sacrifice Creed of Nurture. If you do, reveal a creature from your hand and choose a creature in play. For the remainder of the turn, the chosen creature gains the text box of the revealed creature.
var CreedOfNurture = card.New(
	"Creed of Nurture",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 386),
	card.WithTraits(card.Traits.Power),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
