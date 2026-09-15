//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// BlastFromThePast
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Exalt a friendly creature. Archive a Saurian creature from your discard pile. Deal damage equal to the archived creature's power to an enemy creature.
var BlastFromThePast = set.New(
	"Blast from the Past",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "200"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
