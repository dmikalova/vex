package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Aemberheart
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Choose a friendly creature. Exalt and ward the chosen creature. Fully heal the chosen creature.
var Aemberheart = set.New(
	"Aemberheart",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "142"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ChooseCreatureThen{
			Target: card.Target.FriendlyCreature,
			Then: card.Sequence{Effects: []card.Effect{
				card.Exalt{
					Target: card.Target.TheChosenCreature,
					Amount: 1,
				},
				card.Ward{Target: card.Target.TheChosenCreature},
				card.Heal{
					Fully:  true,
					Target: card.Target.TheChosenCreature,
				},
			}},
		}),
)
