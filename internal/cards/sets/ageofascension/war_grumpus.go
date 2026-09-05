package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// War Grumpus
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Beast
//
//	Fight/Reap: Ready and fight with a neighboring Giant trait creature.
var WarGrumpus = card.New(
	"War Grumpus",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 52),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithFightOrReap(card.OnChooseCreature{
		Target: card.Target.Creature.Neighboring().WithTrait(card.Traits.Giant),
		Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
	}),
)
