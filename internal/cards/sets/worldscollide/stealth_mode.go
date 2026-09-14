package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Stealth Mode
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Until the end of your next turn, players cannot play Tactics.
var StealthMode = set.New(
	"Stealth Mode",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "306"),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.PlayersCannotPlay{
		Type:     card.Type.Tactic,
		Duration: card.Duration.EndOfPlayerNextTurn,
	}),
)
