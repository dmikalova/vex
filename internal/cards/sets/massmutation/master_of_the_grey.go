package massmutation

import "github.com/dmikalova/vex/internal/card"

// Master of the Grey
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  1
//	Traits: Human • Monk
//
//	Your opponent cannot resolve bonus icons on cards they play.
var MasterOfTheGrey = set.New(
	"Master of the Grey",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "169"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Monk),
	card.WithRestrictions(card.Restrictions{BonusIcons: card.Opponent}),
)
