//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TirelessCrocag
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Giant
//
//	Tireless Crocag cannot reap.
//	You may use Tireless Crocag as if it belonged to the active house.
//	If your opponent has no creatures in play, destroy Tireless Crocag.
var TirelessCrocag = card.New(
	"Tireless Crocag",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 47),
	card.WithPower(7),
	card.WithTraits("Giant"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
