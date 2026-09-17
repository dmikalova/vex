package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Ardent Hero
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Taunt.
//	Ardent Hero cannot be dealt damage by Mutant creatures or creatures with power 5 or higher.
var ArdentHero = set.New(
	"Ardent Hero",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "126"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithCannotBeDealtDamageBy(card.DamageSource{
		Trait:    card.Traits.Mutant,
		MinPower: 5,
	}),
)
