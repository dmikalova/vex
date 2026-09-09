package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Philophosaurus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Dinosaur • Philosopher
//
//	Reap: You may look at the top 3 cards of your deck, archive 1, put 1 into your hand, and discard 1.
var Philophosaurus = card.New(
	"Philophosaurus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "207"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Philosopher),
	card.WithAbility(
		card.Trigger.Reap, card.May{Do: card.LookAtTopSort{}}),
)
