//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// OldBoomy
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Goblin • Scientist
//
//	Reap: Reveal cards from the top of your deck until you reveal a Brobnar card or choose to stop. Deal 2D to Old Boomy if a Brobnar card was revealed. Archive each card revealed this way.
var OldBoomy = card.New(
	"Old Boomy",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 45),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
