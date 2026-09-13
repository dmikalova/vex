package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Narp
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  8
//	Armor:  1
//	Traits: Giant
//
//	Each neighboring Creature cannot reap.
var Narp = card.New(
	"Narp",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "12"),
	card.WithPower(8),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Giant),
	card.WithConstant(card.ConstantAbility{
		Target:         card.Target.EachCreature.Neighboring(),
		CannotBeUsedTo: card.UseKinds(card.UseKind.Reap),
	}),
)
