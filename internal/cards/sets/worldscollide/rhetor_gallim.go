package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Rhetor Gallim
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Dinosaur • Philosopher
//
//	Play: Keys cost +3 Æmber during your opponent's next turn.
//	Reap: You may exalt Rhetor Gallim. Keys cost +3 Æmber during your opponent's next turn.
var RhetorGallim = set.New(
	"Rhetor Gallim",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "192"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Philosopher),
	card.WithAbility(
		card.Trigger.Play, card.RaiseKeyCost{
			Player:   card.Opponent,
			Amount:   3,
			Duration: card.Duration.OpponentNextTurn,
		}),
	card.WithAbility(
		card.Trigger.Reap, card.May{Do: card.Sequence{Effects: []card.Effect{
			card.Exalt{
				Target: card.Target.This,
				Amount: 1,
			},
			card.RaiseKeyCost{
				Player:   card.Opponent,
				Amount:   3,
				Duration: card.Duration.OpponentNextTurn,
			},
		}}}),
)
