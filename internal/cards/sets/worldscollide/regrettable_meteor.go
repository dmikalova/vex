package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Regrettable Meteor
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Destroy each Dinosaur creature and each creature with power 6 or higher.
var RegrettableMeteor = card.New(
	"Regrettable Meteor",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 208),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.EachCreature.WithTrait(card.Traits.Dinosaur)},
			card.Destroy{Target: card.Target.EachCreature.PowerAtLeast(6)},
		}}),
)
