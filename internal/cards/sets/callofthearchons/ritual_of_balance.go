package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ritual of Balance
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Power
//
//	Action: If your opponent has 6 Æmber or more, steal 1 Æmber.
var RitualOfBalance = card.New(
	"Ritual of Balance",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 342),
	card.WithTraits(card.Traits.Power),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.OpponentAember{Is: card.AtLeast, Amount: 6},
			Then: card.StealAember{Amount: 1},
		}),
)
