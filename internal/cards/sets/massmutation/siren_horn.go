package massmutation

import "github.com/dmikalova/vex/internal/card"

// Siren Horn
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Before Fight: Move 1 Æmber from this creature to the creature it fights."
var SirenHorn = set.New(
	"Siren Horn",
	card.House.Saurian,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "212"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.BeforeFight,
			Effect: card.MoveAember{
				Amount: 1,
				From:   card.Target.This,
				Onto:   card.Target.CreatureFought,
			},
		}},
	}),
)
