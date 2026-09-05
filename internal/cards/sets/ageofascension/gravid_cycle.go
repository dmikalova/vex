package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Gravid Cycle
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Omega.
//	Play: Put a card from your discard pile into your hand.
var GravidCycle = card.New(
	"Gravid Cycle",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 354),
	card.WithAemberBonus(1),
	card.WithKeywords(card.Keyword.Omega),
	card.WithAbility(
		card.Trigger.Play, card.PutFromDiscard{Destination: card.To.Hand}),
)
