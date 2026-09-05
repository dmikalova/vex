package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Sack of Coins
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Choose a creature - for each Æmber in your pool, deal 1 damage to the chosen creature.
var SackOfCoins = card.New(
	"Sack of Coins",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 312),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.DealDamage{
				Amount: 1,
				Per:    card.AemberInPool{Player: card.Controller},
				Target: card.Target.TheChosenCreature,
			},
		}),
)
