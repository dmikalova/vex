package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Break-key
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If your opponent has more forged keys than you, unforge one of your opponent's keys. Your opponent gains 6 Æmber.
var Breakkey = set.New(
	"Break-key",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "019"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.HasMoreForgedKeys{Player: card.Opponent},
			Then: card.Sentences{Effects: []card.Effect{
				card.UnforgeKey{Player: card.Opponent},
				card.GainAember{Player: card.Opponent, Amount: 6},
			}},
		}),
)
