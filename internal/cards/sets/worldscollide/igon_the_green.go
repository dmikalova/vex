//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// IgonTheGreen
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	Destroyed: Purge Igon the Green. Return an Igon the Terrible from your discard pile to your hand.
var IgonTheGreen = card.New(
	"Igon the Green",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "39"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
