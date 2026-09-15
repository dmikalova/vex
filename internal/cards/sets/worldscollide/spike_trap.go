package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Spike Trap
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Spike Trap -> deal 3 damage to each flank creature.
var SpikeTrap = set.New(
	"Spike Trap",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "261"),
	card.WithBonus(card.Bonus.Aember),
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
