//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// ShadowSelf
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  9
//	Traits: Specter
//
//	Shadow Self deals no damage when fighting.
//	Damage dealt to non-Specter neighbors is dealt to Shadow Self instead.
var ShadowSelf = card.New(
	"Shadow Self",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 310),
	card.WithPower(9),
	card.WithTraits("Specter"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
