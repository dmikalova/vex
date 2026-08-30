//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// RedPlanetRayGun
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Reap: Choose a creature. Deal 1 Damage to that creature for each Mars creature in play."
var RedPlanetRayGun = card.New(
	"Red Planet Ray Gun",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 211),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
