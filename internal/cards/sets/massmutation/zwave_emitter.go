//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ZWaveEmitter
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	At the start of your turn, ward this creature.
var ZWaveEmitter = set.New(
	"Z-Wave Emitter",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Special,
	card.Provenance(card.MM, "356"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
