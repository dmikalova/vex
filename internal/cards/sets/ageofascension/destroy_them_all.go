//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// DestroyThemAll
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy an artifact, a creature,
//	and an upgrade.
var DestroyThemAll = card.New(
	"Destroy Them All!",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 179),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
