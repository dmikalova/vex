//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TitanMechanic
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Cyborg • Scientist
//
//	While Titan Mechanic is on a flank, each key costs -1 Aember.
var TitanMechanic = card.New(
	"Titan Mechanic",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 154),
	card.WithPower(6),
	card.WithTraits("Cyborg", "Scientist"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
