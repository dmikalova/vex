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
//	Reap: Discard cards from the top of your deck until you discard a Brobnar card or choose to stop -> deal 2 damage to Old Boomy. Archive each card discarded this way.
var OldBoomy = set.New(
	"Old Boomy",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "45"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Reap, card.Sentences{Effects: []card.Effect{
			card.Then{
				First: card.DiscardTopOfDeckUntil{
					House:   card.House.Self,
					MayStop: true,
				},
				Result: card.DealDamage{Target: card.Target.This, Amount: 2},
			},
			card.ArchiveDiscardedThisWay{},
		}}),
)
