package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Cowfyne's Brew
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Give a creature two +1 power counters.
var CowfynesBrew = card.New(
	"Cowfyne's Brew",
	card.House.Brobnar,
	card.Type.Tactic,
	// TODO(variant): rarity relabelled from Variant to Rare — handle manually
	card.Rarity.Rare,
	// TODO(duplicate): mechanically identical to Alaka's Brew — fold/handle manually
	card.Provenance(card.WC, 63),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.Creature,
			Amount: 2,
		}),
)
