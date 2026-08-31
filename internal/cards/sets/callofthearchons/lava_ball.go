package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lava Ball
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Deal 4 damage to a creature that is not on a flank and 2 damage to each of its neighbors.
var LavaBall = card.New(
	"Lava Ball",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 9),
	card.WithAbility(
		card.Trigger.Play, card.SplashDamage{
			Amount: 4,
			Splash: 2,
		}),
)
