package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Berinon
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  2
//	Traits: Spirit • Knight
//
//	After a creature enters play, if it is a Mutant creature, enrage Berinon.
//	Reap: Berinon captures 2 Æmber from your opponent.
var Berinon = set.New(
	"Berinon",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "146"),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Spirit, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.AfterCreatureEnters, card.Conditional{
			Cond: card.ItIsOfTrait{Trait: card.Traits.Mutant},
			Then: card.Enrage{Target: card.Target.This},
		}),
	card.WithAbility(
		card.Trigger.Reap, card.CaptureAember{
			Amount: 2,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
