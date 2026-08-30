package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dimension Door
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For the remainder of the turn, instead of gaining Æmber from reaping, steal the same amount.
var DimensionDoor = card.New(
	"Dimension Door",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 108),
	card.WithAbility(
		card.Trigger.Play, card.Instead{
			Of:   card.Event.ReapAember,
			With: card.Steal,
		}),
)
