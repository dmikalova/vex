package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Kartanoo
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast
//
//	Reap: Use an artifact.
var Kartanoo = set.New(
	"Kartanoo",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "348"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.Reap, card.Use{
			Max:    1,
			Target: card.Target.EachArtifact,
		}),
)
