package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pile of Skulls
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Each time an enemy creature is destroyed during your turn, a friendly creature captures 1 Æmber from your opponent.
var PileOfSkulls = card.New(
	"Pile of Skulls",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 25),
	card.WithTraits("Location"),
	card.WithAbility(
		card.Trigger.AfterEnemyCreatureDestroyed, card.CaptureAember{
			Amount: 1,
			Target: card.Target.FriendlyCreature,
			Source: card.Opponent,
		}),
)
