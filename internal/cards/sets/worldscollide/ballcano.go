package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ballcano
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Deal 4 damage to each creature. Gain 2 chains.
var Ballcano = card.New(
	"Ballcano",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 3),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.DealDamage{
				Target: card.Target.EachCreature,
				Amount: 4,
			},
			card.GainChains{
				Player: card.Controller,
				Amount: 2,
			},
		}}),
)
