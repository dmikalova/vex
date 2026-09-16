package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Eunoia
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Beast • Cat
//
//	After a creature is destroyed in a fight with Eunoia, gain 1 Æmber, and heal 2 damage from Eunoia.
var Eunoia = set.New(
	"Eunoia",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "400"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Beast, card.Traits.Cat),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.Sequence{Effects: []card.Effect{
			card.GainAember{Player: card.Controller, Amount: 1},
			card.Heal{Amount: 2, Target: card.Target.This},
		}}),
)
