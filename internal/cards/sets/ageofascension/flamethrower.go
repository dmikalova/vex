package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Flamethrower
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Weapon
//
//	Action: Deal 1 damage to a creature that is not on a flank and 1 damage to each of its neighbors.
var Flamethrower = card.New(
	"Flamethrower",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 21),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.DealDamage{Spread: card.CreatureAndNeighbors{
			Amount:     1,
			Splash:     1,
			NotOnFlank: true,
		}}),
)
