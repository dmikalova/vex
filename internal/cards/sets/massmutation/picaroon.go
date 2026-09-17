package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Picaroon
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  0
//	Traits: Mutant • Changeling
//
//	Deploy.
//	X is the combined power of Picaroon's non-Changeling neighbors.
var Picaroon = set.New(
	"Picaroon",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "028"),
	card.WithPower(0),
	card.WithTraits(card.Traits.Mutant, card.Traits.Changeling),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithPowerX(card.CombinedPowerOfNeighborsWithout{
		Without: card.Traits.Changeling,
	}),
)
