//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// CrazyKillingMachine
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Weapon
//
//	Action: Discard the top card of each player's deck. For each of those cards, destroy a creature or artifact of that card's house, if able. If 2 cards are not destroyed as a result of this, destroy Crazy Killing Machine.
var CrazyKillingMachine = card.New(
	"Crazy Killing Machine",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 141),
	card.WithTraits(card.Traits.Weapon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
