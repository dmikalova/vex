package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// First Officer Frane
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Play/Fight/Reap: A friendly creature captures 1 Æmber from your opponent.
var FirstOfficerFrane = card.New(
	"First Officer Frane",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 298),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	card.WithPlayFightReap(card.CaptureAember{
		Amount: 1,
		Target: card.Target.FriendlyCreature,
		Source: card.Opponent,
	}),
)
