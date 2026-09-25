package worldscollide

import (
	"slices"

	"github.com/dmikalova/vex/internal/card"
)

// upgradeOrRobot matches the cards Chief Engineer Walls retrieves — any Upgrade,
// or any card with the Robot trait — so a deck that runs Walls is guaranteed a
// couple of them to pull back from the discard pile.
func upgradeOrRobot(d card.Definition) bool {
	if d.Type == card.Type.Upgrade {
		return true
	}
	return slices.Contains(d.Traits, card.Traits.Robot)
}

// Chief Engineer Walls
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human
//
//	Elusive.
//	Play/Fight/Reap: You may put an upgrade or Robot card from your discard pile into your hand.
var ChiefEngineerWalls = set.New(
	"Chief Engineer Walls",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "293"),
	card.InCluster(card.Pulled(wallsBlasterCluster, 1, 1.25)),
	card.PullsMatching("Walls' Upgrades and Robots", 2, 2, upgradeOrRobot),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(card.Trigger.PlayFightReap, card.May{
		Do: card.PutCard{Zones: []card.Zone{card.Discard},
			Selection: card.Chosen{
				Type: card.Type.Upgrade,
				Or:   []card.Filter{{Trait: card.Traits.Robot}},
			},
			Destination: card.To.Hand,
		},
	}),
)
