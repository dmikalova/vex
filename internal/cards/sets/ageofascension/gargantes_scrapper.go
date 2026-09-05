//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// GargantesScrapper
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Giant
//
//	Alpha. (You can only play this card before doing anything else this step.)
//	Play: If you have 3A or more,
//	deal 3D to an enemy creature.
var GargantesScrapper = card.New(
	"Gargantes Scrapper",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 23),
	card.WithPower(3),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
