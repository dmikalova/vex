//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ZParticleTracker
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gains, "Fight: Search your deck for an upgrade and put it into your hand. Shuffle your deck."
var ZParticleTracker = set.New(
	"Z-Particle Tracker",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Special,
	card.Provenance(card.MM, "354"),
	card.WithBonus(card.Bonus.Aember),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
