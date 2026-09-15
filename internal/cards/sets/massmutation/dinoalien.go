//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DinoAlien
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Mutant • Alien
//
//	Play: You may exalt Dino-Alien. If you do, deal 3D to a creature.
//	Fight: Look at the top 3 cards of your deck. Put 1 into your hand and 1 on the bottom of your deck.
var DinoAlien = set.New(
	"Dino-Alien",
	card.House.Staralliance,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "357"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Alien),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
