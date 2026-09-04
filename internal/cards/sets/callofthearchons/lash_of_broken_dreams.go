package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lash of Broken Dreams
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Keys cost +3 Æmber during your opponent's next turn.
var LashOfBrokenDreams = card.New(
	"Lash of Broken Dreams",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, 75),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.RaiseKeyCost{
			Player:   card.Opponent,
			Amount:   3,
			Duration: card.Duration.NextTurn,
		}),
)
