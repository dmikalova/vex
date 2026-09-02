package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Skippy Timehog
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Play: Your opponent cannot use any cards during their next turn.
var SkippyTimehog = card.New(
	"Skippy Timehog",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 152),
	card.WithPower(1),
	card.WithTraits("Mutant"),
	card.WithAbility(
		card.Trigger.Play, card.CannotUse{
			Player:   card.Opponent,
			Duration: card.Duration.NextTurn,
		}),
)
