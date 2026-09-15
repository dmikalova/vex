//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// EnsignElSamra
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Mutant
//
//	Enhance Draw Draw Draw.
//	Action: Reveal a card from your hand. Resolve its bonus icons as if you had played it.
var EnsignElSamra = set.New(
	"Ensign El-Samra",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "340"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
