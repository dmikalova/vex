package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Molephin
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Hazardous 3.
//	After Æmber is stolen from you, for each Æmber stolen, deal 1 damage to each enemy Creature.
var Molephin = set.New(
	"Molephin",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "360"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithHazardous(3),
	card.WithAbility(
		card.Trigger.AfterAemberStolenFromYou, card.DealDamage{
			Amount: 1,
			Per:    card.AemberStolenThisEvent{},
			Target: card.Target.EachEnemyCreature,
		}),
)
