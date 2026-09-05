package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// The Grey Rider
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Monk
//
//	Deploy.
//	Play/Fight/Reap: You may ready and fight with a neighboring creature.
var TheGreyRider = card.New(
	"The Grey Rider",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 226),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Monk),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithPlayFightReap(card.May{
		Do: card.OnChooseCreature{
			Target: card.Target.Creature.Neighboring(),
			Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
		},
	}),
)
