package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mogghunter's Brew
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//	Æmber:  1
//
//	Play: Give a creature two +1 power counters.
var MogghuntersBrew = card.New(
	"Mogghunter's Brew",
	card.House.Brobnar,
	card.Type.Tactic,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	// TODO(duplicate): mechanically identical to Alaka's Brew — fold/handle manually
	card.Provenance(card.WC, "66"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.Creature,
			Amount: 2,
		}),
)
