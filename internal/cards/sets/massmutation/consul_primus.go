package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Consul Primus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Dinosaur • Politician
//
//	Reap: Move 1 Æmber from a creature to another creature.
//	Enhance Capture.
var ConsulPrimus = set.New(
	"Consul Primus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "187"),
	card.InCluster(card.Pulled(monumentToPrimusCluster, 1, 1)),
	card.WithEnhance(card.Bonus.Capture),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithAbility(
		card.Trigger.Reap, card.MoveAember{
			Amount: 1,
			From:   card.Target.Creature,
			Onto:   card.Target.OtherCreature,
		}),
)
