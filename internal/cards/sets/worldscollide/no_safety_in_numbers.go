package worldscollide

import "github.com/dmikalova/vex/internal/card"

// No Safety in Numbers
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Deal 3 damage to each creature that belongs to a house that has 3 or more creatures in play.
var NoSafetyInNumbers = set.New(
	"No Safety in Numbers",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "257"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 3,
			Target: card.Target.EachCreature.Refine(card.HouseWithAtLeast(3)),
		}),
)
