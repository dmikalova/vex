//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// HelmsmanSpears
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Human
//
//	Fight/Reap: Discard any number of cards from your hand. Draw a card for each card discarded this way.
var HelmsmanSpears = card.New(
	"Helmsman Spears",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 311),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
