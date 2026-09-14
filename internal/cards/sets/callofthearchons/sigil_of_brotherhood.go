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
//	Action: Destroy Sigil of Brotherhood. For the remainder of the turn, you may use friendly Sanctum Creatures.
var SigilOfBrotherhood = set.New(
	"Sigil of Brotherhood",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "236"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Power),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(card.Trigger.Action, card.Sentences{Effects: []card.Effect{
		card.Destroy{Target: card.Target.This},
		card.MayPlayOrUse{Houses: card.GrantHouses.Named(card.House.Self), Grant: card.GrantUse},
	}}),
)
