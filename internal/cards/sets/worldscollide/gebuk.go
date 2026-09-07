//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Gebuk
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Beast
//
//	Destroyed: Discard the top card of your deck. If it is a creature, after Gebuk leaves play, put that creature into play in Gebuk's position in the battleline.
var Gebuk = card.New(
	"Gebuk",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 373),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
