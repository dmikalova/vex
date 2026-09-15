package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Ritual of Tognath
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber Æmber Æmber
//
//	Play: Destroy 2 friendly creatures.
var RitualOfTognath = set.New(
	"Ritual of Tognath",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "044"),
	card.WithBonus(card.Bonus.Aember, card.Bonus.Aember, card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DestroyChosen{
			Target: card.Target.EachFriendlyCreature,
			Amount: 2,
		}),
)
