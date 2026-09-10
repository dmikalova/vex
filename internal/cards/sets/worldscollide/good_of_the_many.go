package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Good of the Many
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy each creature that does not share a trait with another creature in its controller's battleline.
var GoodOfTheMany = card.New(
	"Good of the Many",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "220"),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.Refine(card.WithoutSharedTrait()),
		}),
)
