package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Kompsos Haruspex
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Priest
//
//	Each friendly Creature's play effect is a play/reap effect.
var KompsosHaruspex = set.New(
	"Kompsos Haruspex",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "224"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Priest),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachFriendlyCreature,
		Morphs: []card.TriggerMorph{{
			From: card.Trigger.Play,
			Onto: card.Trigger.Reap,
		}},
	}),
)
