//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LesserOxtet
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Demon
//
//	Elusive.
//	Play: Purge each card in your hand.
//	Reap: Keys cost +3A during your opponent's next turn.
var LesserOxtet = card.New(
	"Lesser Oxtet",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 109),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
