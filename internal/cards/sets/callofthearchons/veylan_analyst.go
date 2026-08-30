//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// VeylanAnalyst
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	Each time you use an artifact, gain 1 Aember.
var VeylanAnalyst = card.New(
	"Veylan Analyst",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 156),
	card.WithPower(2),
	card.WithTraits("Cyborg", "Scientist"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
