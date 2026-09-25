package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Chota Hazri
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human • Witch
//
//	Play: Lose 1 Æmber. Forge a key at current cost -> purge Chota Hazri.
var ChotaHazri = set.New(
	"Chota Hazri",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "349"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.LoseAember{
				Player: card.Controller,
				Amount: 1,
			},
			card.ForgeKey{},
		}}),
)
