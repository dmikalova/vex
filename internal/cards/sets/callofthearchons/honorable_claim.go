package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Honorable Claim
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Each friendly Knight trait creature captures 1 Æmber from your opponent.
var HonorableClaim = card.New(
	"Honorable Claim",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 219),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EachFriendlyCreature.WithTrait("Knight"),
			Source: card.Opponent,
		}),
)
