package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// J43G3R V
//
//	House:  Star Alliance
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  8
//	Armor:  2
//	Traits: Robot
//
//	Reap: Reap with 2 non-Star Alliance creatures, one at a time.
//	Fight: Fight with 2 non-Star Alliance creatures, one at a time.
var J43G3RV = set.Gigantic(
	"J43G3R V",
	card.House.StarAlliance,
	card.Rarity.Rare,
	card.Provenance(card.MoMu, "331"),
	card.WithPower(8),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.Reap, card.Use{
			Max:    2,
			Verb:   card.ReapVerb{},
			Target: card.Target.EachFriendlyCreature.House(card.Houses.Except(card.House.Self)),
		}),
	card.WithAbility(
		card.Trigger.Fight, card.Use{
			Max:    2,
			Verb:   card.FightVerb{},
			Target: card.Target.EachFriendlyCreature.House(card.Houses.Except(card.House.Self)),
		}),
)
