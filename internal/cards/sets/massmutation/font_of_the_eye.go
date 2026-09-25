package massmutation

import "github.com/dmikalova/vex/internal/card"

// Font of the Eye
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Versatile.
//	Action: If an enemy creature has been destroyed this turn, a friendly creature captures 1 Æmber from your opponent.
var FontOfTheEye = set.New(
	"Font of the Eye",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.MM, "134"),
	card.WithTraits(card.Traits.Location),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.CreatureDestroyedThisTurn{Player: card.Opponent},
			Then: card.CaptureAember{
				Amount: 1,
				Target: card.Target.FriendlyCreature,
				Source: card.Opponent,
			},
		}),
)
