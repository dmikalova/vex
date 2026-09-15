package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lifeward
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Lifeward. Your opponent cannot play creatures during their next turn.
var Lifeward = set.New(
	"Lifeward",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "77"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Power),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(card.Trigger.Action, card.Sentences{Effects: []card.Effect{
		card.Destroy{Target: card.Target.This},
		card.CannotPlay{
			Player:   card.Opponent,
			Type:     card.Type.Creature,
			Duration: card.Duration.OpponentNextTurn,
		},
	}}),
)
