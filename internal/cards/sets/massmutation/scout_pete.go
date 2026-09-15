//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ScoutPete
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Alien
//
//	Play/Fight/Reap: Look at the top card of your deck. You may discard that card.
var ScoutPete = set.New(
	"Scout Pete",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "311"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Alien),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
