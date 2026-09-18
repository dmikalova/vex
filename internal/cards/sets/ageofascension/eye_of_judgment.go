package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Eye of Judgment
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Purge a creature from a discard pile.
var EyeOfJudgment = set.New(
	"Eye of Judgment",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "253"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.PurgeCard{
			Zone:      card.Discard,
			Player:    card.ChosenPlayer,
			Selection: card.Chosen{Type: card.Type.Creature},
		}),
)
