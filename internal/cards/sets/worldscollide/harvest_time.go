//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// HarvestTime
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a trait. Purge each card with that trait. Each player gains 1A for each card they controlled that was purged this way.
var HarvestTime = card.New(
	"Harvest Time",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "106"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
