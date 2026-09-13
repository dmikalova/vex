package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Silent Dagger
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains, "Reap: Deal 4 damage to a flank Creature."
var SilentDagger = card.New(
	"Silent Dagger",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "318"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.Reap, Effect: card.DealDamage{
				Amount: 4,
				Target: card.Target.Creature.OnFlank(),
			}},
		},
	}),
)
