package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Key Hammer
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If your opponent forged a key during their previous turn, unforge one of your opponent's keys -> your opponent gains 6 Æmber.
var KeyHammer = set.New(
	"Key Hammer",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "66"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.ForgedKey{
				Player:   card.Opponent,
				Previous: true,
			},
			Then: card.Then{
				First: card.UnforgeKey{Player: card.Opponent},
				Result: card.GainAember{
					Player: card.Opponent,
					Amount: 6,
				},
			},
		}),
)
