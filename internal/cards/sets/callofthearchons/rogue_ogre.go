//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// RogueOgre
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant • Mutant
//
//	At the end of your turn, if you played exactly one card this turn, Rogue Ogre heals 2 damage and captures 1 Aember.
var RogueOgre = card.New(
	"Rogue Ogre",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 45),
	card.WithPower(6),
	card.WithTraits("Giant", "Mutant"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
