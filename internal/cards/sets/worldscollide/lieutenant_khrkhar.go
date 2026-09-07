//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LieutenantKhrkhar
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Alien • Handuhan
//
//	Taunt. (This creature's neighbors cannot be attacked unless they have taunt.)
//	Hazardous 3. (Before this creature is attacked, deal 3D to the attacking enemy.)
var LieutenantKhrkhar = card.New(
	"Lieutenant Khrkhar",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 299),
	card.WithPower(5),
	card.WithTraits(card.Traits.Alien, card.Traits.Handuhan),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
