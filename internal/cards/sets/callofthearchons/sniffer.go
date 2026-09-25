package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Sniffer
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Ally
//
//	Action: For the remainder of the turn, each creature loses elusive.
var Sniffer = set.New(
	"Sniffer",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "188"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Ally),
	card.WithAbility(
		card.Trigger.Action, card.LoseKeyword{Keyword: card.Keyword.Elusive}),
)
