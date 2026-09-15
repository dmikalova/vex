//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Turnkey
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Demon
//
//	Play: Unforge an opponent's key. If you do, when Turnkey leaves play, your opponent forges a key at no cost.
var Turnkey = set.New(
	"Turnkey",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "051"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
