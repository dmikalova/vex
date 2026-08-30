//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MartiansMakeBadAllies
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Reveal your hand. Purge each revealed non-Mars creature and gain 1 Aember for each card purged this way.
var MartiansMakeBadAllies = card.New(
	"Martians Make Bad Allies",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 168),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
