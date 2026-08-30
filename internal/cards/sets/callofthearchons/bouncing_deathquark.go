package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bouncing Deathquark
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy an enemy creature and a friendly creature -> if there is a friendly creature in play, you may repeat this effect.
var BouncingDeathquark = card.New(
	"Bouncing Deathquark",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 107),
	card.WithAbility(
		card.Trigger.Play, card.MayRepeat{
			Cond: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
			},
			Do: card.Sequence{
				Effects: []card.Effect{
					card.Destroy{Target: card.Target.EnemyCreature},
					card.Destroy{Target: card.Target.FriendlyCreature},
				},
			},
		}),
)
