package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Extinction
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature - destroy each creature that shares a trait with it. Gain 1 chain.
var Extinction = card.New(
	"Extinction",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 196),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Sentences{Effects: []card.Effect{
				card.Destroy{Target: card.Target.EachCreature.SharingTrait()},
				card.GainChains{Amount: 1},
			}},
		}),
)
