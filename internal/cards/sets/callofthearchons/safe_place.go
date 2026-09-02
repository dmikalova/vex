package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Safe Place
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	You may spend Æmber on Safe Place when forging keys.
//	Action: Move 1 Æmber from your pool to Safe Place.
var SafePlace = card.New("Safe Place",
	card.House.Shadows, card.Type.Artifact, card.Rarity.Rare,
	card.Provenance(card.CotA, 289),
	card.WithAemberBonus(1),
	card.WithTraits("Location"),
	card.WithSpendableAember(),
	card.WithAbility(
		card.Trigger.Action, card.MoveAemberFromPool{
			Amount: 1,
			Target: card.Target.This,
		}),
)
