package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Discombobulator
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains, "Your Æmber cannot be stolen."
var Discombobulator = card.New(
	"Discombobulator",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "149"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{AemberCannotBeStolen: true}),
)
