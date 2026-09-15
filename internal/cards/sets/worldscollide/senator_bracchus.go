package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Senator Bracchus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Politician
//
//	You may spend Æmber on friendly creatures as if it were in your pool.
//	Fight/Reap: Exalt Senator Bracchus.
var SenatorBracchus = set.New(
	"Senator Bracchus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "229"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithConstant(card.ConstantAbility{
		Target:            card.Target.EachFriendlyCreature,
		SpendAemberOnCard: true,
	}),
	card.WithAbility(card.Trigger.FightReap, card.Exalt{
		Target: card.Target.This,
		Amount: 1,
	}),
)
