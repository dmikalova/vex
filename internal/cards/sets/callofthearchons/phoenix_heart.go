package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Phoenix Heart
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains, "Destroyed: Put this creature into its owner's hand, and deal 3 damage to each creature."
var PhoenixHeart = card.New(
	"Phoenix Heart",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 51),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.Destroyed, Effect: card.Sequence{Effects: []card.Effect{
				card.PutFromPlay{
					Target:      card.Target.This,
					Destination: card.To.Hand,
				},
				card.DealDamage{
					Amount: 3,
					Target: card.Target.EachCreature,
				},
			}}},
		},
	}),
)
