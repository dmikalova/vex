package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Stealth Mode
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Your opponent cannot play Tactics during their next turn.
var StealthMode = card.New(
	"Stealth Mode",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "306"),
	// TODO(duplicate): mechanically identical to Scrambler Storm (Logos) — fold/handle manually.
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.CannotPlay{
		Player:   card.Opponent,
		Type:     card.Type.Tactic,
		Duration: card.Duration.NextTurn,
	}),
)
