package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Hidden Stash
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Archive a card from your hand.
var HiddenStash = card.New(
	"Hidden Stash",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 286),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ArchiveFromHand{Amount: 1}),
)
