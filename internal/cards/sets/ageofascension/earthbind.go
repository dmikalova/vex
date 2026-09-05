//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Earthbind
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature cannot be used unless its controller has discarded a card this turn.
var Earthbind = card.New(
	"Earthbind",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 352),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
