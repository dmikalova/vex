package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Q-Mechs
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Robot
//
//	Play: Draw a card.
//	Destroyed: Archive Q-Mechs.
var QMechs = set.New(
	"Q-Mechs",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "078"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.Play, card.Draw{Amount: 1}),
	card.WithAbility(
		card.Trigger.Destroyed, card.ArchiveSource{}),
)
