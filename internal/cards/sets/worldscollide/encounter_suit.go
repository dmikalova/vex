//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// EncounterSuit
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Rare
//
//	After an action card is played but before it resolves, ward this creature.
var EncounterSuit = card.New(
	"Encounter Suit",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "330"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
