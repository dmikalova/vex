//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// OneLastJob
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Purge each friendly Shadows creature. Steal 1 Aember for each creature purged this way.
var OneLastJob = card.New(
	"One Last Job",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 277),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
