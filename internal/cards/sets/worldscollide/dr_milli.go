//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DrMilli
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Scientist
//
//	Play: For each creature your opponent controls in excess of you, not counting Dr. Milli, archive a card.
var DrMilli = card.New(
	"Dr. Milli",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 150),
	card.WithPower(2),
	card.WithTraits(card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
