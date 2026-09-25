package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Techivore Pulpate
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Jelly
//
//	After a player chooses an active house, destroy each artifact of that house.
var TechivorePulpate = set.New(
	"Techivore Pulpate",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "341"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Jelly),
	card.WithEachPlayerAbility(
		card.Trigger.AfterChooseHouse,
		card.Destroy{Target: card.Target.EachArtifact.House(card.Houses.Active)},
	),
)
