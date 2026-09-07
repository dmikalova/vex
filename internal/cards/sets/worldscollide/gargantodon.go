//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Gargantodon
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  16
//	Traits: Beast
//
//	Gargantodon enters play stunned.
//	Gargantodon only deals 4D when fighting.
//	Each A that would be stolen is captured by a creature controlled by the active player instead.
var Gargantodon = card.New(
	"Gargantodon",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 203),
	card.WithPower(16),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
