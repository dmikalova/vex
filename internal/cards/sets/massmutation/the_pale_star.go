package massmutation

import "github.com/dmikalova/vex/internal/card"

// The Pale Star
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Versatile.
//	Action: Destroy The Pale Star. For the remainder of the turn, each creature is considered to have 1 power and 0 armor.
var ThePaleStar = set.New(
	"The Pale Star",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "049"),
	card.WithTraits(card.Traits.Power),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.OverrideStats{
				Power:    1,
				HasPower: true,
				Armor:    0,
				HasArmor: true,
				Duration: card.Duration.RemainderOfPlayerTurn,
			},
		}}),
)
