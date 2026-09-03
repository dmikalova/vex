package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Reverse Time
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Swap your deck and your discard pile, then shuffle your deck.
var ReverseTime = card.New(
	"Reverse Time",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 121),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.SwapDeckAndDiscard{}),
)
