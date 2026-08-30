//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TendrilsOfPain
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 1 Damage to each creature. Deal an additional 3 Damage to each creature if your opponent forged a key on their previous turn.
var TendrilsOfPain = card.New(
	"Tendrils of Pain",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 64),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
