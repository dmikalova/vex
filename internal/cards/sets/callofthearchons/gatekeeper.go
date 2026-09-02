package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gatekeeper
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Armor:  1
//	Traits: Knight • Spirit
//
//	Play: If your opponent has 7 Æmber or more, Gatekeeper captures all but 5 Æmber from your opponent.
var Gatekeeper = card.New(
	"Gatekeeper",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 260),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits("Knight", "Spirit"),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAember{
				Is:     card.AtLeast,
				Amount: 7,
			},
			Then: card.CaptureAember{
				By:     card.AllBut(5),
				Target: card.Target.This,
				Source: card.Opponent,
			},
		}),
)
