package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// No Safety in Numbers
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 3 damage to each Creature that belongs to a house that has 3 or more Creatures in play.
var NoSafetyInNumbers = card.New(
	"No Safety in Numbers",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "257"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 3,
			Target: card.Target.EachCreature.Refine(card.HouseWithAtLeast(3)),
		}),
)
