//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DaemoAlien
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Alien
//
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
//	Destroyed: Steal 1A.
var DaemoAlien = set.New(
	"Daemo-Alien",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "306"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Alien),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
