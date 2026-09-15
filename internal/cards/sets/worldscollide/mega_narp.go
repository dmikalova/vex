package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Narp
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  10
//	Armor:  1
//	Traits: Giant
//
//	Each neighboring creature cannot reap.
var MegaNarp = set.New(
	"Mega Narp",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "60"),
	card.InCluster(card.Pulled(narpsBrewCluster, 1, 1.25)),
	card.WithPower(10),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Giant),
	card.WithConstant(card.ConstantAbility{
		Target:         card.Target.EachCreature.Neighboring(),
		CannotBeUsedTo: card.UseKinds(card.UseKind.Reap),
	}),
)
