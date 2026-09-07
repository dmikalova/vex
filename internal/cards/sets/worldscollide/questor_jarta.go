//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// QuestorJarta
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Dinosaur • Politician
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Reap: You may exalt Questor Jarta. If you do, gain 1A.
var QuestorJarta = card.New(
	"Questor Jarta",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 191),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
