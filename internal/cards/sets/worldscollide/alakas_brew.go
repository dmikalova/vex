package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Alaka's Brew
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Give a creature 2 +1 power counters.
var AlakasBrew = card.New(
	"Alaka's Brew",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 2),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.Creature,
			Amount: 2,
		}),
)
