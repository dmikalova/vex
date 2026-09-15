package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Call to Action
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Ready each friendly Knight creature.
var CallToAction = set.New(
	"Call to Action",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "160"),
	card.WithAbility(
		card.Trigger.Play, card.Ready{
			Target: card.Target.EachFriendlyCreature.WithTrait(card.Traits.Knight),
		}),
)
