package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Punctuated Equilibrium
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Each player discards their hand, then refills their hand as if it were the end of their turn.
var PunctuatedEquilibrium = card.New(
	"Punctuated Equilibrium",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 363),
	card.WithAbility(
		card.Trigger.Play, card.EachPlayerDiscardsAndRefillsHand{}),
)
