//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SciOfficerQincan
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Alien • Proximan • Scientist
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	After a player chooses an active house which matches no cards in play, steal 1A.
var SciOfficerQincan = card.New(
	"Sci. Officer Qincan",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "304"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Alien, card.Traits.Proximan, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
