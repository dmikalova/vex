//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// OperationsOfficerYshi
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Spirit
//
//	Taunt.
//	Each of Operations Officer Yshi's neighbors gains, "Fight/Reap: Capture 1A."
var OperationsOfficerYshi = card.New(
	"Operations Officer Yshi",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 334),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Spirit),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
