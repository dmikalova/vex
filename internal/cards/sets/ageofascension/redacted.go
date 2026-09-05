//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// REDACTED
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: [redacted]
//
//	After you choose Logos as your active house, place 1A from the common supply on [REDACTED]. When there are 4 or more A on [REDACTED], sacrifice it and forge a key at no cost.
var REDACTED = card.New(
	"[REDACTED]",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 139),
	card.WithTraits(card.Traits.Redacted),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
