package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Tempting Offer
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Put an enemy creature into its owner's hand -> your opponent gains 1 Æmber.
//	Enhance Capture.
var TemptingOffer = set.New(
	"Tempting Offer",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "259"),
	card.WithBonus(card.Bonus.Aember),
	card.WithEnhance(card.Bonus.Capture),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.PutFromPlay{
				Target:      card.Target.EnemyCreature,
				Destination: card.To.Hand,
			},
			Result: card.GainAember{
				Player: card.Opponent,
				Amount: 1,
			},
		}),
)
