//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// IgonTheTerrible
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: FIXED
//	Power:  8
//	Traits: Giant
//
//	Play: If Igon the Green has not been purged, destroy Igon the Terrible.
//	Fight: Steal 1A.
var IgonTheTerrible = card.New(
	"Igon the Terrible",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 53),
	card.WithPower(8),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
