//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Senator Shrix
var SenatorShrix = card.New(
	"Senator Shrix",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 193),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithSpendableAember(),
	card.WithPlayReap(card.May{Do: card.Exalt{
		Target: card.Target.This,
		Amount: 1,
	}}),
)
