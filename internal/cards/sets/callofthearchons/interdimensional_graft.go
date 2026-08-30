package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Interdimensional Graft
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If an opponent forges a key on their next turn, they must give you their remaining Æmber.
var InterdimensionalGraft = card.New(
	"Interdimensional Graft",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 112),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.GiveRemainingAemberAfterOpponentForgeKey{}),
)
