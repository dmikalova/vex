package massmutation

import "github.com/dmikalova/vex/internal/card"

// Access Denied
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature cannot reap.
var AccessDenied = set.New(
	"Access Denied",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.MM, "302"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		CannotBeUsedTo: card.UseKinds(card.UseKind.Reap),
	}),
)
