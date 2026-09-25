package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Sack of Coins
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Choose a creature. For each Æmber in your pool, deal 1 damage to the chosen creature.
var SackOfCoins = set.New(
	"Sack of Coins",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "312"),
	card.WithBonus(card.Bonus.Aember),
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
