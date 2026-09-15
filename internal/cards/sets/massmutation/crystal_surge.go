package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Crystal Surge
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Exalt each Mutant creature.
var CrystalSurge = set.New(
	"Crystal Surge",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "217"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Target: card.Target.EachCreature.WithTrait(card.Traits.Mutant),
			Amount: 1,
		}),
)
