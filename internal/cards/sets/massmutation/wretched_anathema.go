package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Wretched Anathema
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  10
//	Traits: Demon
//
//	While there are no other friendly creatures in play, Wretched Anathema gains, "Action: Gain 4 Æmber."
//	Play/Reap: Destroy 2 other creatures.
var WretchedAnathema = set.Gigantic(
	"Wretched Anathema",
	card.House.Dis,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "023"),
	card.WithPower(10),
	card.WithTraits(card.Traits.Demon),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.This,
		WhileCondition: card.CardsInPlay{
			Player: card.Controller,
			Type:   card.Type.Creature,
			Other:  true,
			None:   true,
		},
		Granted: []card.Ability{{
			Trigger: card.Trigger.Action,
			Effect:  card.GainAember{Player: card.Controller, Amount: 4},
		}},
	}),
	card.WithAbility(
		card.Trigger.PlayReap, card.DestroyChosen{
			Target: card.Target.OtherCreature,
			Amount: 2,
		}),
)
