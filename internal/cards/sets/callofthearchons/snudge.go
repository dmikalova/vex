//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Snudge
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Demon
//
//	Fight/Reap: Return an artifact or flank creature to its owner's hand.
var Snudge = card.New(
	"Snudge",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 97),
	card.WithPower(4),
	card.WithTraits("Demon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
