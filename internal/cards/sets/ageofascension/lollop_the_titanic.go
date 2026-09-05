//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// LollopTheTitanic
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  11
//	Traits: Giant • Location
//
//	Lollop the Titanic deals no damage when attacked.
var LollopTheTitanic = card.New(
	"Lollop the Titanic",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 14),
	card.WithPower(11),
	card.WithTraits(card.Traits.Giant, card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
