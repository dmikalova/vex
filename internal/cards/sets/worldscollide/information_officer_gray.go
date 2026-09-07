//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// InformationOfficerGray
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
//	Play/Fight/Reap: You may reveal a non-Star Alliance card from your hand. If you do, archive it.
var InformationOfficerGray = card.New(
	"Information Officer Gray",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 312),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
