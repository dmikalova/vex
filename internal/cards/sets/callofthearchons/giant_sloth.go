//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// GiantSloth
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Beast
//
//	You cannot use this card unless you have discarded an Untamed card from your hand this turn.
//	Action: Gain 3 Aember.
var GiantSloth = card.New(
	"Giant Sloth",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 354),
	card.WithPower(6),
	card.WithTraits("Beast"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
