//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// NavigatorAli
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: Look at the top 3 cards of your deck and put them back in any order.
var NavigatorAli = card.New(
	"Navigator Ali",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 314),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
