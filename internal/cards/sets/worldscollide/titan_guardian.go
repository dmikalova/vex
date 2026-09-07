//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TitanGuardian
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Beast • Cyborg
//
//	Taunt. (This creature's neighbors cannot be attacked unless they have taunt.)
//	Destroyed: If Titan Guardian is not on a flank, draw 2 cards.
var TitanGuardian = card.New(
	"Titan Guardian",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 141),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Beast, card.Traits.Cyborg),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
