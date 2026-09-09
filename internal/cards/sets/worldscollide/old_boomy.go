package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Old Boomy
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Goblin • Scientist
//
//	Reap: Reveal cards from the top of your deck until you reveal a Brobnar card or choose to stop, archiving each card revealed this way -> deal 2 damage to Old Boomy.
var OldBoomy = card.New(
	"Old Boomy",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "45"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Reap, card.Then{
			First:  card.RevealDeckUntilHouse{House: card.House.Self},
			Result: card.DealDamage{Target: card.Target.This, Amount: 2},
		}),
)
