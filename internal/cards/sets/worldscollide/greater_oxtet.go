//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GreaterOxtet
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Taunt.
//	At the end of your "ready cards" step, purge a card from your hand. If you do, give Greater Oxtet two +1 power counters.
var GreaterOxtet = card.New(
	"Greater Oxtet",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "105"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
