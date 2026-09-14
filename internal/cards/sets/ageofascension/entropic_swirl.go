package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Entropic Swirl
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a Creature - for each trait that Creature has, deal 2 damage to the chosen Creature, and gain 1 Æmber.
var EntropicSwirl = set.New(
	"Entropic Swirl",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "143"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.ForEach{
				Times: card.TraitsOfChosen{},
				Do: card.Sequence{Effects: []card.Effect{
					card.DealDamage{
						Amount: 2,
						Target: card.Target.TheChosenCreature,
					},
					card.GainAember{
						Player: card.Controller,
						Amount: 1,
					},
				}},
			},
		}),
)
