package massmutation

import "github.com/dmikalova/vex/internal/card"

// Mushroom with a View
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Location
//
//	Versatile.
//	Action: Heal 1 damage from each friendly creature.
var MushroomWithAView = set.New(
	"Mushroom with a View",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "386"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Heal{
			Amount: 1,
			Target: card.Target.EachFriendlyCreature,
		}),
)
