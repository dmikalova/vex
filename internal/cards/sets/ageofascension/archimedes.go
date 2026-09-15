package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Archimedes
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg • Beast
//
//	Elusive.
//	Each neighboring creature gains, "Destroyed: Archive this creature from play."
var Archimedes = set.New(
	"Archimedes",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "108"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Beast),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature.Neighboring(),
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect:  card.ArchiveFromPlay{Target: card.Target.This},
		}},
	}),
)
