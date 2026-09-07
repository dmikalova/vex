package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Quantum Fingertrap
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Action: Swap the positions of two creatures in a battleline.
var QuantumFingertrap = card.New(
	"Quantum Fingertrap",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "133"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.SwapChosen{}),
)
