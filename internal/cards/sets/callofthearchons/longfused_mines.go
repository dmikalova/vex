package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Longfused Mines
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	Versatile.
//	Action: Destroy Longfused Mines, and deal 3 damage to each enemy creature that is not on a flank.
var LongfusedMines = card.New(
	"Longfused Mines",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 287),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Weapon),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.DealDamage{
				Amount: 3,
				Target: card.Target.EachEnemyCreature.NotOnFlank(),
			},
		}}),
)
