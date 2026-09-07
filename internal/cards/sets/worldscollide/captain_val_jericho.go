//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CaptainValJericho
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Human • Leader
//
//	During your turn, if Captain Val Jericho is in the center of your battleline, you may play one card that is not of the active house.
var CaptainValJericho = card.New(
	"Captain Val Jericho",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "326"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Leader),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
