//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Noname
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Demon
//
//	Noname gets +1 power for each purged card.
//	Play/Fight/Reap: Purge a card in a discard pile.
var Noname = card.New(
	"Noname",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 112),
	card.WithPower(1),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
