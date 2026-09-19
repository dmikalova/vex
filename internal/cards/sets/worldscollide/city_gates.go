package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// City Gates
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: A friendly creature captures 1 Æmber from your opponent. If it is a Dinosaur creature, it captures 1 Æmber from your opponent.
var CityGates = set.New(
	"City Gates",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "216"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.CaptureAember{
				Amount: 1,
				Target: card.Target.FriendlyCreature,
				Source: card.Opponent,
			},
			card.Conditional{
				Cond: card.ItIsOfTrait{Trait: card.Traits.Dinosaur},
				Then: card.CaptureAember{
					Amount: 1,
					Target: card.Target.Triggering,
					Source: card.Opponent,
				},
			},
		}}),
)
