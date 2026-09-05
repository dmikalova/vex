//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// HeistNight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha. (You can only play this card before doing anything else this step.)
//	Play: Steal 1A for each friendly Thief creature.
var HeistNight = card.New(
	"Heist Night",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 303),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
