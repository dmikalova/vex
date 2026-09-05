package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Skullion
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Armor:  2
//	Traits: Demon
//
//	Play: Destroy a friendly creature.
var Skullion = card.New(
	"Skullion",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 82),
	card.WithPower(7),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{Target: card.Target.FriendlyCreature}),
)
