package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Red Planet Ray Gun
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Reap: For each Mars creature in play, deal 1 damage to a creature."
var RedPlanetRayGun = card.New(
	"Red Planet Ray Gun",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 211),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect: card.DealDamage{
				Amount: 1,
				Target: card.Target.Creature,
				Per:    card.InPlay{Player: card.EachPlayer, Type: card.Type.Creature, House: card.House.Mars},
			},
		}},
	}),
)
