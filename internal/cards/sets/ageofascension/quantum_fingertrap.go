package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Quantum Fingertrap
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Swap the positions of two creatures in a battleline.
var QuantumFingertrap = set.New(
	"Quantum Fingertrap",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "133"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.SwapChosen{}),
)
