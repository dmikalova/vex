//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Manchego
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Thief
//
//	Play: If you have 5 or fewer cards in your deck, steal 2A.
//	Fight/Reap: You may shuffle Manchego into your deck.
var Manchego = card.New(
	"Manchego",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "275"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Thief),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
