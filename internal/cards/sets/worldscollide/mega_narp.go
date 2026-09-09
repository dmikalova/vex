package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Narp
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  10
//	Armor:  1
//	Traits: Giant
//
//	Each neighboring creature cannot reap.
var MegaNarp = card.New(
	"Mega Narp",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "60"),
	card.WithPower(10),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Giant),
	card.WithConstant(card.ConstantAbility{
		Target:         card.Target.EachCreature.Neighboring(),
		CannotBeUsedTo: card.UseKinds(card.UseKind.Reap),
	}),
)
