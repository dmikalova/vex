//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hyde
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Human • Scientist
//
//	Reap: Draw a card. If you control Velum, draw 2 cards instead.
//	Destroyed: Archive Velum from your discard pile. If you do, archive Hyde.
var Hyde = card.New(
	"Hyde",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "167"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
