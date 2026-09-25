package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Chant of Hubris
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Move 1 Æmber from a creature to another creature.
var ChantOfHubris = set.New(
	"Chant of Hubris",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "184"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.MoveAember{
			Amount: 1,
			From:   card.Target.Creature,
			Onto:   card.Target.OtherCreature,
		}),
)
