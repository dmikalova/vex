package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Stealth Mode
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Until the end of your next turn, players cannot play tactics.
var StealthMode = set.New(
	"Stealth Mode",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "306"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.PlayersCannotPlay{
		Type:     card.Type.Tactic,
		Duration: card.Duration.EndOfPlayerNextTurn,
	}),
)
