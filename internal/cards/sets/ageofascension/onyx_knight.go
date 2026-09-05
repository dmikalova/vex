package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Onyx Knight
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon • Knight
//
//	Play: Destroy each creature with odd power.
var OnyxKnight = card.New(
	"Onyx Knight",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 95),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.OddPower(),
		}),
)
