package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ANT1-10NY
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Robot
//
//	At the end of your turn, move 1 Æmber from ANT1-10NY to your opponent's pool.
//	Play: ANT1-10NY captures all your opponent's Æmber.
var ANT110NY = set.New(
	"ANT1-10NY",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "318"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.MoveAember{
			Amount: 1,
			From:   card.Target.This,
			To:     card.Opponent,
		}),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			All:    true,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
