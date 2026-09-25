package massmutation

import "github.com/dmikalova/vex/internal/card"

// Essence Scale
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Choose a friendly creature. Destroy the chosen creature. Ready and use a friendly creature of that card's house.
var EssenceScale = set.New(
	"Essence Scale",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "021"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ChooseCreatureThen{
			Target: card.Target.FriendlyCreature,
			Then: card.Sequence{Effects: []card.Effect{
				card.Destroy{Target: card.Target.TheChosenCreature},
				card.OnChooseCreature{
					Target: card.Target.FriendlyCreature.House(card.Houses.Contextual),
					Verbs: []card.CreatureVerb{
						card.ReadyVerb{},
						card.UseVerb{},
					},
				},
			}},
		}),
)
