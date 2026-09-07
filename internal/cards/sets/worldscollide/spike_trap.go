package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Spike Trap
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Spike Trap -> deal 3 damage to each flank creature.
var SpikeTrap = card.New(
	"Spike Trap",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "261"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Weapon),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.Destroy{Target: card.Target.This},
			Result: card.DealDamage{
				Amount: 3,
				Target: card.Target.EachCreature.OnFlank(),
			},
		}),
)
