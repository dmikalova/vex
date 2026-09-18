package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gongoozle
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Deal 3 damage to a creature. If it is not destroyed, its owner discards a random card from their hand.
var Gongoozle = set.New(
	"Gongoozle",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "60"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.DealDamage{
		Amount: 3,
		After:  card.IfSurvives,
		Target: card.Target.Creature,
		Then: card.DiscardCard{
			Player:    card.ItsOwner,
			Zones:     []card.Zone{card.Hand},
			Selection: card.Random{Count: 1},
		},
	}),
)
