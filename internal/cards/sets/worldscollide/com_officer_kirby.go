package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Com. Officer Kirby
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: Play a non-Star Alliance Artifact, Upgrade, or Tactic.
var ComOfficerKirby = set.New(
	"Com. Officer Kirby",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "295"),
	card.InCluster(card.Pulled(kirbysBlasterCluster, 1, 1.25)),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(card.Trigger.PlayFightReap, card.PlayFrom{
		From:  card.Hand,
		House: card.Houses.Except(card.House.Self),
		Types: card.Types.Of(card.Type.Artifact, card.Type.Upgrade, card.Type.Tactic),
	}),
)
