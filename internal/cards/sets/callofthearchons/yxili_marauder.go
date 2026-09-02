package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Yxili Marauder
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Martian • Soldier
//
//	Yxili Marauder gains +1 power for each Æmber on it.
//	Play: For each friendly ready Mars creature, Yxili Marauder captures 1 Æmber from your opponent.
var YxiliMarauder = card.New(
	"Yxili Marauder",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 203),
	card.WithPower(2),
	card.WithTraits("Martian", "Soldier"),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.This,
		PowerBonus: 1,
		Per:        card.AemberOnThis{},
	}),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.House.Self,
				Ready:  true,
			},
		}),
)
