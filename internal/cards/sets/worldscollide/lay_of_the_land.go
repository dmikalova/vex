//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LayOfTheLand
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Look at the top 3 cards of your deck and put them back in any order. Draw a card.
var LayOfTheLand = card.New(
	"Lay of the Land",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 313),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
