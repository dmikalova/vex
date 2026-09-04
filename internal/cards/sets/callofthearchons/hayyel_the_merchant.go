package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hayyel the Merchant
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Merchant
//
//	After you play an artifact, gain 1 Æmber.
var HayyelTheMerchant = card.New(
	"Hayyel the Merchant",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 245),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Merchant),
	card.WithAbility(card.Trigger.AfterCardPlayed, card.Conditional{
		Cond: card.ItIs{Type: card.Type.Artifact},
		Then: card.GainAember{
			Player: card.Controller,
			Amount: 1,
		},
	}),
)
