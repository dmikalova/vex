package massmutation

import "github.com/dmikalova/vex/internal/card"

// Patronage
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Move half the Æmber from a creature to your pool, rounding up. Move all Æmber from the chosen creature to your opponent's pool.
var Patronage = set.New(
	"Patronage",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "227"),
	card.WithAbility(card.Trigger.Play, card.Sequence{Effects: []card.Effect{
		card.MoveAember{
			From:     card.Target.Creature,
			Fraction: card.HalfRoundedUp,
			To:       card.Controller,
			Bind:     true,
		},
		card.MoveAember{
			From: card.Target.TheChosenCreature,
			All:  true,
			To:   card.Opponent,
		},
	}}),
)
