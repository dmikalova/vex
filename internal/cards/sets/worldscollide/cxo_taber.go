//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CXOTaber
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Alien • Krxix
//
//	Fight/Reap: You may play or use one non-Star Alliance card this turn.
var CXOTaber = card.New(
	"CXO Taber",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 309),
	card.WithPower(3),
	card.WithTraits(card.Traits.Alien, card.Traits.Krxix),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
