package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Special Delivery
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Versatile.
//	Action: Deal 3 damage to a flank creature. If this damage destroys that creature, purge it.
var SpecialDelivery = card.New(
	"Special Delivery",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 292),
	card.WithAemberBonus(1),
	card.WithTraits("Item"),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.DamageIfDestroyed{
			Amount: 3,
			Target: card.Target.Creature.OnFlank(),
			Then:   card.PurgeCreature{Target: card.Target.Triggering},
		}),
)
