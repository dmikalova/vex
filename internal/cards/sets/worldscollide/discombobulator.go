package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Discombobulator
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Your Æmber cannot be stolen."
var Discombobulator = set.New(
	"Discombobulator",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "149"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{AemberCannotBeStolen: card.AlwaysMet{}}),
)
