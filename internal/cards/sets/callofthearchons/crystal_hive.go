package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Crystal Hive
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Action: For the remainder of the turn, after a creature reaps, gain 1 Æmber.
var CrystalHive = card.New(
	"Crystal Hive",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 182),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.ForRemainderOfTurn{
			On: card.Event.Reap,
			Do: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
