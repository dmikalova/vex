package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Special Delivery
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Action: Deal 3 damage to a flank creature. If this damage destroys that creature, purge it.
var SpecialDelivery = set.New(
	"Special Delivery",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "292"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.DealDamage{
			Amount: 3,
			After:  card.IfDestroyed,
			Target: card.Target.Creature.OnFlank(),
			Then:   card.PurgeCreature{Target: card.Target.Triggering},
		}),
)
