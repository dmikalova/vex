//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// StealerOfSouls
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Demon
//
//	After an enemy creature is destroyed fighting Stealer of Souls, purge that creature and gain 1 Aember.
var StealerOfSouls = card.New(
	"Stealer of Souls",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 98),
	card.WithPower(6),
	card.WithTraits("Demon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
