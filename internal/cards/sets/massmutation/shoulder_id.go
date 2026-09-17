package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Shoulder Id
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Specter
//
//	Taunt.
//	When Shoulder Id would deal damage, steal 1 Æmber instead.
//	Shoulder Id cannot fight.
var ShoulderID = set.New(
	"Shoulder Id",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "257"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Specter),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithCannotBeUsedTo(card.UseKind.Fight),
	card.WithStealsInsteadOfDamageWhenAttacked(1),
)
