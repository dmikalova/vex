//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SaurusRex
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Dinosaur • Leader
//
//	Fight/Reap: If Saurus Rex is in the center of your battleline, you may exalt it. If you do, search your deck for a Saurian card, reveal it, and add it to your hand. Then, shuffle your deck.
var SaurusRex = card.New(
	"Saurus Rex",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 227),
	card.WithPower(6),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Leader),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
