package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Clear Mind
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Unstun each friendly creature.
var ClearMind = set.New(
	"Clear Mind",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "216"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Unstun{Target: card.Target.EachFriendlyCreature}),
)
