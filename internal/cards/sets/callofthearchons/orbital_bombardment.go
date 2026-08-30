//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// OrbitalBombardment
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Reveal any number of Mars cards from your hand. For each card revealed this way, deal 2 Damage to a creature. (You may choose a different creature each time.)
var OrbitalBombardment = card.New(
	"Orbital Bombardment",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 172),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
