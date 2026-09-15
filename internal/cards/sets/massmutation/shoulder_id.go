//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// ShoulderId
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Specter
//
//	Taunt.
//	Shoulder Id cannot fight.
//	When Shoulder Id would deal damage, steal 1A instead.
var ShoulderId = set.New(
	"Shoulder Id",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "257"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Specter),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
