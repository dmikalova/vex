package massmutation

import "github.com/dmikalova/vex/internal/card"

// Mole
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Your opponent may spend Æmber on this creature as if it were in their pool."
var Mole = set.New(
	"Mole",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.MM, "287"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{SpendAemberOnCard: card.SpendScope.Opponent}),
)
