package worldscollide

import "github.com/dmikalova/vex/internal/card"

// The Golden Spiral
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Action: Choose a friendly creature. Exalt the chosen creature. Ready and use the chosen creature.
var TheGoldenSpiral = set.New(
	"The Golden Spiral",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.WC, "194"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.ChooseCreatureThen{
			Target: card.Target.FriendlyCreature,
			Then: card.Sequence{Effects: []card.Effect{
				card.Exalt{
					Target: card.Target.TheChosenCreature,
					Amount: 1,
				},
				card.OnChooseCreature{
					Target: card.Target.TheChosenCreature,
					Verbs: []card.CreatureVerb{
						card.ReadyVerb{},
						card.UseVerb{},
					},
				},
			}},
		}),
)
