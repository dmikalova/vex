package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Ortannu the Chained
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Demon
//
//	Reap: Put each Ortannu's Binding from your discard pile into your hand. For each card put into your hand this way, deal 2 damage to a creature and 2 damage to each of its neighbors.
var OrtannuTheChained = card.New(
	"Ortannu the Chained",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 97),
	card.WithPower(7),
	card.WithTraits(card.Traits.Demon),
	card.Connects(
		card.Pull(OrtannusBinding, 2),
	),
	card.WithAbility(
		card.Trigger.Reap, card.Sentences{
			Effects: []card.Effect{
				card.PutFromDiscard{
					Name:        OrtannusBinding.Name,
					All:         true,
					Destination: card.To.Hand,
				},
				card.Repeat{
					Times: card.CardsReturnedThisWay{},
					Do: card.DealDamage{Spread: card.CreatureAndNeighbors{
						Amount: 2,
						Splash: 2,
					}},
				},
			},
		}),
)
