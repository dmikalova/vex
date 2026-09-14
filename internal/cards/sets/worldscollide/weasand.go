package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Weasand
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast • Thief
//
//	Deploy, Elusive.
//	If Weasand is on a flank, destroy Weasand.
//	After a player forges a key, gain 2 Æmber.
var Weasand = set.New(
	"Weasand",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "285"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Deploy, card.Keyword.Elusive),
	card.WithDestroyedWhen(card.OnFlank{}),
	card.WithAbility(
		card.Trigger.AfterPlayerForgesKey, card.GainAember{
			Player: card.Controller,
			Amount: 2,
		}),
)
