package massmutation

import "github.com/dmikalova/vex/internal/card"

// Gizelhart's Wrath
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy each Mutant creature.
var GizelhartsWrath = set.New(
	"Gizelhart's Wrath",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "163"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.WithTrait(card.Traits.Mutant),
		}),
)
