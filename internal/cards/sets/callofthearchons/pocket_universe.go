package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pocket Universe
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	You may spend Æmber on Pocket Universe when forging keys.
//	Action: Move 1 Æmber from your pool to Pocket Universe.
var PocketUniverse = card.New("Pocket Universe",
	card.House.Logos, card.Type.Artifact, card.Rarity.Rare,
	card.Provenance(card.CotA, 131),
	card.WithTraits("Item"),
	card.WithSpendableAember(),
	card.WithAbility(
		card.Trigger.Action, card.MoveAemberFromPool{
			Amount: 1,
			Target: card.Target.This,
		}),
)
