//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// KelifiDragon
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  12
//	Traits: Dragon
//
//	Kelifi Dragon cannot be played unless you have 7 Aember or more.
//	Fight/Reap: Gain 1 Aember. Deal 5 Damage to a creature.
var KelifiDragon = card.New(
	"Kelifi Dragon",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 37),
	card.WithPower(12),
	card.WithTraits("Dragon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
