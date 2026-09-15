//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LCdrTrigon
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant
//
//	Reap: Discard the top card of your deck. Resolve that card's bonus icons as if you had played it.
var LCdrTrigon = set.New(
	"LCdr. Trigon",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "324"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
