package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Terms of Redress
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: A friendly creature captures 2 Æmber from your opponent.
var TermsOfRedress = set.New(
	"Terms of Redress",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "227"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 2,
			Target: card.Target.FriendlyCreature,
			Source: card.Opponent,
		}),
)
