package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Bumblebird
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Beast • Insect
//
//	Alpha.
//	Play: Give each other friendly Untamed creature two +1 power counters.
var Bumblebird = set.New(
	"Bumblebird",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "336"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast, card.Traits.Insect),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.EachOtherFriendlyCreature.House(card.Houses.Named(card.House.Self)),
			Amount: 2,
		}),
)
