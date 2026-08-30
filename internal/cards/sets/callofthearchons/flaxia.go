package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Flaxia
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Faerie
//
//	Play: If you control more creatures than your opponent, gain 2 Æmber.
var Flaxia = card.New(
	"Flaxia",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 352),
	card.WithPower(4),
	card.WithTraits("Faerie"),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.ControlsMoreCreatures{},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 2,
			},
		}),
)
