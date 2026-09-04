package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bumpsy
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Play: Your opponent loses 1 Æmber.
var Bumpsy = card.New(
	"Bumpsy",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 30),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Play, card.LoseAember{
			Player: card.Opponent,
			Amount: 1,
		}),
)
