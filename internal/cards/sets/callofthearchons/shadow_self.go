package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Shadow Self
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  9
//	Traits: Specter
//
//	Shadow Self deals no damage when fighting.
//	Damage dealt to each neighboring non-Specter trait creature is dealt to Shadow Self instead.
var ShadowSelf = card.New(
	"Shadow Self",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 310),
	card.WithPower(9),
	card.WithTraits(card.Traits.Specter),
	card.WithAttackDamage(card.AttackDamage{Fixed: true}),
	card.WithTakesDamageFor(
		card.Target.EachCreature.Neighboring().ExceptTrait(card.Traits.Specter)),
)
