package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Spoils of Battle
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: A friendly creature captures 1 Æmber from your opponent. Each creature with Æmber on it captures 1 Æmber from its opponent.
var SpoilsOfBattle = set.New(
	"Spoils of Battle",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "197"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.CaptureAember{
				Amount: 1,
				Target: card.Target.FriendlyCreature,
				Source: card.Opponent,
			},
			card.CaptureAember{
				Amount: 1,
				Target: card.Target.EachCreature.WithAember(),
				Source: card.ItsOpponent,
			},
		}}),
)
