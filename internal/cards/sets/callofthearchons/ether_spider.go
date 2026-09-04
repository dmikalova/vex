package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ether Spider
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Traits: Beast
//
//	Ether Spider deals no damage when fighting.
//	If Æmber would be added to your opponent's pool, instead Ether Spider captures it.
var EtherSpider = card.New(
	"Ether Spider",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 192),
	card.WithPower(7),
	card.WithTraits(card.Traits.Beast),
	card.WithAttackDamage(card.AttackDamage{
		Amount: 0,
		Fixed:  true,
	}),
	card.WithReplaces(card.Instead{
		Of:     card.Event.AemberAddedToPool,
		Player: card.Opponent,
		With:   card.Capture,
	}),
)
