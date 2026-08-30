//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// NumquidTheFair
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	Play: Destroy an enemy creature. Repeat this card's effect if your opponent still controls more creatures than you.
var NumquidTheFair = card.New(
	"Numquid the Fair",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 253),
	card.WithPower(3),
	card.WithTraits("Human"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
