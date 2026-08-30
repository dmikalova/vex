//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// JehuTheBureaucrat
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	After you choose Sanctum as your active house, gain 2 Aember.
var JehuTheBureaucrat = card.New(
	"Jehu the Bureaucrat",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 250),
	card.WithPower(3),
	card.WithTraits("Human"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
