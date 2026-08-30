package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lomir Flamefist
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Giant
//
//	Play: If your opponent has 7 Æmber or more, your opponent loses 2 Æmber.
var LomirFlamefist = card.New(
	"Lomir Flamefist",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 40),
	card.WithPower(5),
	card.WithTraits("Giant"),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAemberAtLeast{Amount: 7},
			Then: card.LoseAember{
				Player: card.Opponent,
				Amount: 2,
			},
		}),
)
