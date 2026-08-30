package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Clear Mind
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Unstun each friendly creature.
var ClearMind = card.New(
	"Clear Mind",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 216),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Unstun{Target: card.Target.EachFriendlyCreature}),
)
