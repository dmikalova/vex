package massmutation

import "github.com/dmikalova/vex/internal/card"

// Detention Coil
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature cannot fight.
var DetentionCoil = set.New(
	"Detention Coil",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "321"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		CannotBeUsedTo: card.UseKinds(card.UseKind.Fight),
	}),
)
