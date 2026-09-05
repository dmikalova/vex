//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// StormCrawler
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Robot
//
//	Storm Crawler only deals 1D when fighting.
//	After an enemy creature reaps, stun it.
var StormCrawler = card.New(
	"Storm Crawler",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 189),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Robot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
