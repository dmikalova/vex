package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Techivore Pulpate
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Jelly
//
//	After a player chooses an active house, destroy each artifact of that house.
var TechivorePulpate = card.New(
	"Techivore Pulpate",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "341"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Jelly),
	card.WithAbility(
		card.Trigger.AfterAnyPlayerChoosesHouse,
		card.Destroy{Target: card.Target.EachArtifact.OfActiveHouse()},
	),
)
