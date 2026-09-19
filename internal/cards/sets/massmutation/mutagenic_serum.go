package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mutagenic Serum
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Mutagenic Serum. For the remainder of the turn, you may use friendly Mutant creatures.
var MutagenicSerum = set.New(
	"Mutagenic Serum",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "091"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(card.Trigger.Action, card.Sequence{Effects: []card.Effect{
		card.Destroy{Target: card.Target.This},
		card.MayPlayOrUse{Trait: card.Traits.Mutant, Grant: card.GrantUse},
	}}),
)
