package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Masterplan
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Play: Put a card from your hand facedown under Masterplan.
//	Action: Play the card under Masterplan. Destroy Masterplan.
var Masterplan = set.New(
	"Masterplan",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "288"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Play, card.PutUnderFromHand{FaceDown: true}),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.PlayCardUnder{},
			card.Destroy{Target: card.Target.This},
		}}),
)
