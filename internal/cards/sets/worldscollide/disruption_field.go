//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DisruptionField
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	Your opponent's keys cost +1A for each disruption counter on Disruption Field.
//	This creature gains "Fight/Reap: Put a disruption counter on Disruption Field."
var DisruptionField = card.New(
	"Disruption Field",
	card.House.Staralliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, 328),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
