//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// YxiliMarauder
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Martian • Soldier
//
//	Yxili Marauder gets +1 power for each Aember on it.
//	Play: Capture 1 Aember for each friendly ready Mars creature.
var YxiliMarauder = card.New(
	"Yxili Marauder",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 203),
	card.WithPower(2),
	card.WithTraits("Martian", "Soldier"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
