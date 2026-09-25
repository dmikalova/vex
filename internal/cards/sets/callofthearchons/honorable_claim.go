package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Honorable Claim
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Each friendly Knight creature captures 1 Æmber from your opponent.
var HonorableClaim = set.New(
	"Honorable Claim",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "219"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EachFriendlyCreature.WithTrait(card.Traits.Knight),
			Source: card.Opponent,
		}),
)
