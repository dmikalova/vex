package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sniffer
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Ally
//
//	Action: For the remainder of the turn, each creature loses elusive.
var Sniffer = card.New(
	"Sniffer",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 188),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Ally),
	card.WithAbility(
		card.Trigger.Action, card.LoseKeyword{Keyword: card.Keyword.Elusive}),
)
