//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Aemberlution
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Omega.
//	Play: Destroy each creature. Each player reveals their hand and puts each creature from their hand into play ready.
var Aemberlution = set.New(
	"Aemberlution",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "394"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
