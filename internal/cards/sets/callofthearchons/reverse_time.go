package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Reverse Time
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Swap your deck and your discard pile, then shuffle your deck.
var ReverseTime = set.New(
	"Reverse Time",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "121"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.SwapDeckAndDiscard{}),
)
