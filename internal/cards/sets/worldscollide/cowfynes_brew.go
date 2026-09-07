//go:build todo

// TODO: source rarity is "Variant"; awaiting a real-rarity mapping before implementing.
package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Cowfyne's Brew
var CowfynesBrew = card.New(
	"Cowfyne's Brew",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 63),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.Creature,
			Amount: 2,
		}),
)
