//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// QincansBlaster
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Variant
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may deal 2D to a creature, or attach Qincan's Blaster to Sci. Officer Qincan."
//	After you attach Qincan's Blaster to Sci. Officer Qincan, you may archive a creature in play.
var QincansBlaster = card.New(
	"Qincan's Blaster",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 351),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
