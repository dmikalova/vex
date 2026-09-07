//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TacticalOfficerMoon
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Assault 2. (Before this creature attacks, deal 2D to the attacked enemy.)
//	Play: You may rearrange the creatures in a player's battleline.
var TacticalOfficerMoon = card.New(
	"Tactical Officer Moon",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "320"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
