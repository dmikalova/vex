//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// VezymaThinkdrone
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Scientist
//
//	Reap: You may archive a friendly creature or artifact from play.
var VezymaThinkdrone = card.New(
	"Vezyma Thinkdrone",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 202),
	card.WithPower(3),
	card.WithTraits("Martian", "Scientist"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
