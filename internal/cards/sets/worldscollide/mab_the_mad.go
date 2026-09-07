//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MabTheMad
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Æmber:  1
//	Traits: Faerie
//
//	Reap: Shuffle Mab the Mad into your deck.
var MabTheMad = card.New(
	"Mab the Mad",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 378),
	card.WithPower(2),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Faerie),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
