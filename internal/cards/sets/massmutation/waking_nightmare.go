package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Waking Nightmare
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Keys cost +1 Æmber for each Dis creature in play during your opponent's next turn.
//	Enhance Capture.
var WakingNightmare = set.New(
	"Waking Nightmare",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "017"),
	card.WithBonus(card.Bonus.Aember),
	card.WithEnhance(card.Bonus.Capture),
	card.WithAbility(
		card.Trigger.Play, card.RaiseKeyCost{
			Player:   card.Opponent,
			Amount:   1,
			House:    card.Houses.Named(card.House.Dis),
			Duration: card.Duration.OpponentNextTurn,
		}),
)
