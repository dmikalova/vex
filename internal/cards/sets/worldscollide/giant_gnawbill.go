package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Giant Gnawbill
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Beast
//
//	After a player chooses an active house, that player destroys an artifact of that house.
var GiantGnawbill = set.New(
	"Giant Gnawbill",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "390"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.AfterAnyPlayerChoosesHouse,
		card.ByActivePlayer{
			Do: card.Destroy{Target: card.Target.Artifact.House(card.Houses.Active)},
		},
	),
)
