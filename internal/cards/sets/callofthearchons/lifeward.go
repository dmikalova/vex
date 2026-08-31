package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lifeward
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Lifeward. Your opponent cannot play creatures during their next turn.
var Lifeward = card.New(
	"Lifeward",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 77),
	card.WithTraits("Power"),
	card.WithAemberBonus(1),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(card.Trigger.Action, card.Sequence{Effects: []card.Effect{
		card.Sentence{Effect: card.Destroy{Target: card.Target.This}},
		card.Sentence{Effect: card.CannotPlayNextTurn{
			Player: card.Opponent,
			Type:   card.Type.Creature,
		}},
	}}),
)
