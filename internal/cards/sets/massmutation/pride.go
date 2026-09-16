package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Pride
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Demon • Sin
//
//	Reap: Ward each friendly Sin creature.
var Pride = set.New(
	"Pride",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "060"),
	card.InCluster(sinsCluster),
	card.OneCopyPerDeck(),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	card.WithAbility(
		card.Trigger.Reap,
		card.Ward{Target: card.Target.EachFriendlyCreature.WithTrait(card.Traits.Sin)},
	),
)
