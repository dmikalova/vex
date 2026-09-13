package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Livia the Elder
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Philosopher
//
//	Reap: You may exalt Livia the Elder -> each friendly Creature's fight effects and reap effects are fight/reap effects for the remainder of the turn.
var LiviaTheElder = card.New(
	"Livia the Elder",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "225"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Philosopher),
	card.WithAbility(
		card.Trigger.Reap, card.May{Do: card.Then{
			First: card.Exalt{
				Target: card.Target.This,
				Amount: 1,
			},
			Result: card.FuseTriggersForTurn{
				A: card.Trigger.Fight,
				B: card.Trigger.Reap,
			},
		}}),
)
