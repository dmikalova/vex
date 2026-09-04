package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Grommid
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  10
//	Traits: Beast
//
//	You cannot play creatures.
//	After a creature is destroyed fighting Grommid, your opponent loses 1 Æmber.
var Grommid = card.New(
	"Grommid",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 194),
	card.WithPower(10),
	card.WithTraits(card.Traits.Beast),
	card.WithRestrictions(card.Restrictions{CannotPlay: card.Type.Creature}),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.LoseAember{
			Player: card.Opponent,
			Amount: 1,
		}),
)
