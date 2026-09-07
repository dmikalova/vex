//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// KalochStonefather
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant • Leader
//
//	While Kaloch Stonefather is in the center of your battleline, each friendly creature gains skirmish.
var KalochStonefather = card.New(
	"Kaloch Stonefather",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 41),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant, card.Traits.Leader),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
