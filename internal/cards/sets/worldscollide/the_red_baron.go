//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheRedBaron
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: FIXED
//	Power:  4
//	Armor:  1
//	Traits: Cyborg • Pirate
//
//	While your red key is forged, The Red Baron gains, "Reap: Steal 1A."
//	While your opponent's red key is forged, The Red Baron gains elusive.
var TheRedBaron = card.New(
	"The Red Baron",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 0),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Pirate),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
