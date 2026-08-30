//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// UnguardedCamp
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each creature you have in excess of your opponent, a friendly creature captures 1 Aember. Each creature cannot capture more than 1 Aember this way.
var UnguardedCamp = card.New(
	"Unguarded Camp",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 17),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
