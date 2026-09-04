package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Nepenthe Seed
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Nepenthe Seed, and put a card from your discard pile into your hand.
var NepentheSeed = card.New(
	"Nepenthe Seed",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 341),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.PutFromDiscard{Destination: card.To.Hand},
		}}),
)
