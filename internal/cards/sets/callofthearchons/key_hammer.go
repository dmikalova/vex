package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Key Hammer
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If your opponent forged a key on their previous turn, unforge one of your opponent's keys, and your opponent gains 6 Æmber.
var KeyHammer = card.New(
	"Key Hammer",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 66),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Conditional{
				Cond: card.ForgedKey{
					Player:   card.Opponent,
					Previous: true,
				},
				Then: card.UnforgeKey{Player: card.Opponent},
			},
			card.GainAember{
				Player: card.Opponent,
				Amount: 6,
			},
		}}),
)
