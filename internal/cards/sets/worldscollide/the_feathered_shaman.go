//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheFeatheredShaman
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Human • Witch
//
//	Elusive.
//	Fight/Reap: Ward each of The Feathered Shaman's neighbors.
var TheFeatheredShaman = card.New(
	"The Feathered Shaman",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 383),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
