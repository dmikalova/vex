package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Rotgrub
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Beast
//
//	Play: Your opponent loses 1 Æmber.
//	Reap: Archive Rotgrub.
var Rotgrub = set.New(
	"Rotgrub",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "83"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.Play, card.LoseAember{
			Player: card.Opponent,
			Amount: 1,
		}),
	card.WithAbility(
		card.Trigger.Reap, card.ArchiveSource{}),
)
