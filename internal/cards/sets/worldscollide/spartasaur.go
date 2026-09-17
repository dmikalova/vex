package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Spartasaur
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	After a friendly creature is destroyed, destroy each non-Dinosaur creature.
//	Fight: Gain 2 Æmber.
var Spartasaur = set.New(
	"Spartasaur",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "231"),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.AfterCreatureDestroyed, card.Conditional{
			Cond: card.ItIsFriendly{},
			Then: card.Destroy{
				Target: card.Target.EachCreature.ExceptTrait(card.Traits.Dinosaur),
			},
		}),
	card.WithAbility(
		card.Trigger.Fight, card.GainAember{
			Player: card.Controller,
			Amount: 2,
		}),
)
