package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Silent Dagger
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Reap: Deal 4 damage to a flank creature."
var SilentDagger = set.New(
	"Silent Dagger",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "318"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.Reap, Effect: card.DealDamage{
				Amount: 4,
				Target: card.Target.Creature.OnFlank(),
			}},
		},
	}),
)
