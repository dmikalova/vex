//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// StiltKin
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Goblin
//
//	Skirmish. (When you use this creature to fight, it is dealt no damage in return.)
//	After a Giant creature is played adjacent to Stilt-Kin, ready and fight with Stilt-Kin.
var StiltKin = card.New(
	"Stilt-Kin",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "14"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
