//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MegaShorty
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: FIXED
//	Power:  6
//	Traits: Giant
//
//	Assault 4. (Before this creature attacks, deal 4D to the attacked enemy.)
//	Reap: Enrage Mega Shorty.
var MegaShorty = card.New(
	"Mega Shorty",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 61),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
