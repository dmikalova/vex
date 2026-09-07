//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LiviaTheElder
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Philosopher
//
//	Reap: You may exalt Livia the Elder. If you do, each friendly creature's fight effects and reap effects are fight/reap effects for the remainder of the turn.
var LiviaTheElder = card.New(
	"Livia the Elder",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "225"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Philosopher),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
