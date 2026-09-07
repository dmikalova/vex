package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Gravelguts
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Giant
//
//	After a creature is destroyed in a fight with Gravelguts, give Gravelguts two +1 power counters.
var Gravelguts = card.New(
	"Gravelguts",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "22"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.AddPowerCounter{
			Target: card.Target.This,
			Amount: 2,
		}),
)
