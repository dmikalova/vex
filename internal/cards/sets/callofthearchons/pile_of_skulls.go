package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Pile of Skulls
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	After an enemy creature is destroyed during your turn, a friendly creature captures 1 Æmber from your opponent.
var PileOfSkulls = set.New(
	"Pile of Skulls",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "25"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.AfterCreatureDestroyed, card.Conditional{
			Cond: card.And{Conditions: []card.Condition{
				card.ItIsEnemy{},
				card.ItIsYourTurn{},
			}},
			Then: card.CaptureAember{
				Amount: 1,
				Target: card.Target.FriendlyCreature,
				Source: card.Opponent,
			},
		}),
)
