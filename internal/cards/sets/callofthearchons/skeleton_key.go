package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Skeleton Key
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: A friendly creature captures 1 Æmber from your opponent.
var SkeletonKey = card.New(
	"Skeleton Key",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 291),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.CaptureAember{
			Amount: 1,
			Target: card.Target.FriendlyCreature,
			Source: card.Opponent,
		}),
)
