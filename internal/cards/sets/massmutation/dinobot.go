//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DinoBot
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Mutant • Scientist
//
//	Play: You may exalt Dino-Bot. If you do, deal 3D to a creature.
//	Reap: Discard a card from your hand. If you do, draw a card.
var DinoBot = set.New(
	"Dino-Bot",
	card.House.Logos,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.MM, "119"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
