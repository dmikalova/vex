package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Lilithal
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Fight/Reap: Lilithal captures 1 Æmber from your opponent.
var Lilithal = card.New(
	"Lilithal",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 79),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithFightOrReap(card.CaptureAember{
		Amount: 1,
		Target: card.Target.This,
		Source: card.Opponent,
	}),
)
