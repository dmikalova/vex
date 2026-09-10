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
//	Play/Fight/Reap: You may play a non-Star Alliance artifact, upgrade, or Tactic this turn.
var ComOfficerKirby = card.New(
	"Com. Officer Kirby",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "295"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithPlayFightReap(card.MayPlayOffHouse{
		Except:  card.House.Self,
		NotType: card.Type.Creature,
		Grant:   card.GrantPlay,
		Count:   1,
	}),
)
