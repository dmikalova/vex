//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SicSemperTyrannosaurus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Move each A from the most powerful creature to your pool and destroy that creature.
var SicSemperTyrannosaurus = card.New(
	"Sic Semper Tyrannosaurus",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 209),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
