package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Snarette
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	At the end of your turn, Snarette captures 1 Æmber from your opponent.
//	Action: Move each Æmber on Snarette to the common supply.
var Snarette = set.New(
	"Snarette",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "014"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
	card.WithAbility(
		card.Trigger.Action, card.MoveAemberToSupply{
			All:    true,
			Target: card.Target.This,
		}),
)
