//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheGrimReaper
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: FIXED
//	Power:  4
//	Traits: Robot • Specter
//
//	If you are haunted, The Grim Reaper enters play ready. (You are haunted if there are 10 or more cards in your discard pile.)
//	Reap: Purge an enemy creature and a friendly creature.
var TheGrimReaper = card.New(
	"The Grim Reaper",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 0),
	card.WithPower(4),
	card.WithTraits(card.Traits.Robot, card.Traits.Specter),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
