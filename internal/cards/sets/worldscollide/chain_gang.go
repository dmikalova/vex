//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ChainGang
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	After you play Subtle Chain, ready Chain Gang.
//	Action: Steal 1A. Shuffle a Subtle Chain from your discard pile into your deck.
var ChainGang = card.New(
	"Chain Gang",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "252"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
