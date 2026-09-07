//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Chant of Hubris
var ChantOfHubris = card.New(
	"Chant of Hubris",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, 184),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.MoveAember{
			Amount: 1,
			From:   card.Target.Creature,
			Onto:   card.Target.OtherCreature,
		}),
)
