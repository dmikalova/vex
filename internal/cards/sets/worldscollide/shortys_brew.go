//go:build todo

// TODO: source rarity is "Variant"; awaiting a real-rarity mapping before implementing.
package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Shorty's Brew
var ShortysBrew = card.New(
	"Shorty's Brew",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 68),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.Creature,
			Amount: 2,
		}),
)
