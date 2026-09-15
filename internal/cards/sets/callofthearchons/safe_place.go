package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Safe Place
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Location
//
//	You may spend Æmber on Safe Place when forging keys.
//	Action: Move 1 Æmber from your pool to Safe Place.
var SafePlace = set.New("Safe Place",
	card.House.Shadows, card.Type.Artifact, card.Rarity.Rare,
	card.Provenance(card.CotA, "289"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Location),
	card.WithSpendableAember(),
	card.WithAbility(
		card.Trigger.Action, card.MoveAemberFromPool{
			Amount: 1,
			Target: card.Target.This,
		}),
)
