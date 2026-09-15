//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// CommanderDhrxgar
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant
//
//	After an upgrade is attached to Commander Dhrxgar or one of its neighbors, gain 1A.
var CommanderDhrxgar = set.New(
	"Commander Dhrxgar",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "337"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
