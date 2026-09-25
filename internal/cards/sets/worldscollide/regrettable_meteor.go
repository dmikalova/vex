package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Regrettable Meteor
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Destroy each Dinosaur creature and each creature with power 6 or higher.
var RegrettableMeteor = set.New(
	"Regrettable Meteor",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "208"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.Refine(card.AnyOf(
				card.OfTrait(card.Traits.Dinosaur),
				card.PowerAtLeast(6),
			)),
		}),
)
