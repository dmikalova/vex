package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Primus Unguis
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Each friendly creature gains +2 power for each Æmber on Primus Unguis.
//	Reap: Exalt Primus Unguis.
var PrimusUnguis = card.New(
	"Primus Unguis",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "226"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.EachFriendlyCreature,
		PowerBonus: 2,
		Per:        card.AemberOnThis{},
	}),
	card.WithAbility(
		card.Trigger.Reap, card.Exalt{
			Target: card.Target.This,
			Amount: 1,
		}),
)
