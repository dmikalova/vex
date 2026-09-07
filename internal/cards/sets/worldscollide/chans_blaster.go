//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ChansBlaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Variant
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may deal 2D to a creature, or attach Chan's Blaster to Commander Chan."
//	After you attach Chan's Blaster to Commander Chan, you may use another friendly creature.
var ChansBlaster = card.New(
	"Chan's Blaster",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 345),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
