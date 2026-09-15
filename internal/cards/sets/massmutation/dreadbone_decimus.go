//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DreadboneDecimus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Dinosaur • Assassin
//
//	Play/Fight: You may exalt Dreadbone Decimus. If you do, destroy a creature with lower power than Dreadbone Decimus.
var DreadboneDecimus = set.New(
	"Dreadbone Decimus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "204"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Assassin),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
