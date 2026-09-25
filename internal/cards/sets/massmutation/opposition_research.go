package massmutation

import "github.com/dmikalova/vex/internal/card"

// Opposition Research
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent cannot use creatures to reap during their next turn.
//	Enhance Damage.
var OppositionResearch = set.New(
	"Opposition Research",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "077"),
	card.WithEnhance(card.Bonus.Damage),
	card.WithAbility(
		card.Trigger.Play, card.Restrict{
			Player:   card.Opponent,
			Action:   card.Restricted.Reaping,
			Duration: card.Duration.OpponentNextTurn,
		}),
)
