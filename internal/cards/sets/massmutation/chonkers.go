package massmutation

import "github.com/dmikalova/vex/internal/card"

// Chonkers
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Skirmish.
//	After a creature is destroyed in a fight with Chonkers, give Chonkers +1 power counters equal to the number of +1 power counters on Chonkers.
//	Play: Give Chonkers a +1 power counter.
var Chonkers = set.New(
	"Chonkers",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "396"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.AddPowerCounter{
			Target: card.Target.This,
			Equal:  card.PowerCountersOnThis{},
		}),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.This,
			Amount: 1,
		}),
)
