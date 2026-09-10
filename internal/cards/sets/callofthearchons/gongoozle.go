package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gongoozle
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 3 damage to a creature. If it is not destroyed, its owner discards a random card from their hand.
var Gongoozle = card.New(
	"Gongoozle",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "60"),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.DamageThen{
		Amount: 3,
		After:  card.IfSurvives,
		Target: card.Target.Creature,
		Then:   card.DiscardCard{Player: card.ItsOwner, Zone: card.Hand, Selection: card.Random{}},
	}),
)
