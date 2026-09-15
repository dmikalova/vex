//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// AmbassadorLiu
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant • Politician
//
//	Action: Discard a card from your hand. If it is a Dis or Shadows card, steal 1A. If it is a Logos or Untamed card, gain 2A. If it is a Sanctum or Saurian card, capture 3A.
var AmbassadorLiu = set.New(
	"Ambassador Liu",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "335"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
