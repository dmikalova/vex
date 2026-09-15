//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// TheArchivist
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Cyborg
//
//	If you archive The Archivist, archive it faceup.
//	While The Archivist is in your archives, instead of picking up all of your archives, you may choose to pick up any number of cards in your archives.
var TheArchivist = set.New(
	"The Archivist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "113"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Cyborg),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
