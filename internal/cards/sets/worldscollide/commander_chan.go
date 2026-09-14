package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Commander Chan
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Fight/Reap: Use an other Creature.
var CommanderChan = set.New(
	"Commander Chan",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "296"),
	card.InCluster(card.Pulled(chansBlasterCluster, 1, 1.25)),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(card.Trigger.FightReap, card.Use{
		Max:    1,
		Target: card.Target.EachOtherFriendlyCreature,
	}),
)
