package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Senator Shrix
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Dinosaur • Politician
//
//	You may spend Æmber on Senator Shrix when forging keys.
//	Play/Reap: You may exalt Senator Shrix.
var SenatorShrix = set.New(
	"Senator Shrix",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "193"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithSpendableAember(),
	card.WithAbility(card.Trigger.PlayReap, card.May{Do: card.Exalt{
		Target: card.Target.This,
		Amount: 1,
	}}),
)
