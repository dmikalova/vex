//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TransporterPlatform
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Location
//
//	Action: Return a friendly creature and each upgrade attached to it to your hand.
var TransporterPlatform = card.New(
	"Transporter Platform",
	card.House.Staralliance,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "321"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
