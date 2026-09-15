//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ConsulPrimus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Dinosaur • Politician
//
//	Enhance Capture.
//	Reap: Move 1A from a creature to another creature.
var ConsulPrimus = set.New(
	"Consul Primus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "187"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
