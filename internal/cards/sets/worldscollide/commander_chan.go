//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CommanderChan
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Fight/Reap: Use another friendly creature.
var CommanderChan = card.New(
	"Commander Chan",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 296),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
