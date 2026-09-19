package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Pincerator
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	At the end of each player's turn, deal 1 damage to each flank creature.
var Pincerator = set.New(
	"Pincerator",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "289"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithEachPlayerAbility(
		card.Trigger.EndOfTurn, card.DealDamage{
			Amount: 1,
			Target: card.Target.EachCreature.OnFlank(),
		}),
)
