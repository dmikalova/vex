//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Fidgit
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Faerie • Thief
//
//	Elusive.
//	Reap: Discard a random card from your opponent's archives or the top card of their deck. If that card is an action, play it as if it were yours.
var Fidgit = card.New(
	"Fidgit",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "254"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Faerie, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
