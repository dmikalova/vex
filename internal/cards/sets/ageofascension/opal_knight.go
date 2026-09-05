package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Opal Knight
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Spirit • Knight
//
//	Play: Destroy each creature with even power.
var OpalKnight = card.New(
	"Opal Knight",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 260),
	card.WithPower(5),
	card.WithTraits(card.Traits.Spirit, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.EvenPower(),
		}),
)
