package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sigil of Brotherhood
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Sigil of Brotherhood. For the remainder of the turn, you may use friendly Sanctum creatures.
var SigilOfBrotherhood = card.New(
	"Sigil of Brotherhood",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 236),
	card.WithTraits("Power"),
	card.WithAemberBonus(1),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(card.Trigger.Action, card.Sequence{Effects: []card.Effect{
		card.Sentence{Effect: card.Destroy{Target: card.Target.This}},
		card.Sentence{Effect: card.MayUseFriendlyHouse{House: card.House.Sanctum}},
	}}),
)
