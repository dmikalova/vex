package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Red Baron
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Armor:  1
//	Traits: Cyborg • Pirate
//
//	While your opponent's red key is forged, The Red Baron gains elusive.
//	While your red key is forged, The Red Baron gains, "Reap: Steal 1 Æmber."
var TheRedBaron = card.New(
	"The Red Baron",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "A08"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Pirate),
	card.WithConstant(card.ConstantAbility{
		Target:         card.Target.This,
		WhileCondition: card.KeyColorForged{Player: card.Controller, Color: card.KeyColor.Red},
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect:  card.StealAember{Amount: 1},
		}},
	}),
	card.WithConstant(card.ConstantAbility{
		Target:         card.Target.This,
		WhileCondition: card.KeyColorForged{Player: card.Opponent, Color: card.KeyColor.Red},
		Keywords:       card.Keywords(card.Keyword.Elusive),
	}),
)
