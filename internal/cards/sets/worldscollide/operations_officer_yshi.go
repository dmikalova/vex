package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Operations Officer Yshi
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Spirit
//
//	Taunt.
//	Each neighboring creature gains, "Reap: this creature captures 1 Æmber from your opponent."
//	Each neighboring creature gains, "Fight: this creature captures 1 Æmber from your opponent."
var OperationsOfficerYshi = card.New(
	"Operations Officer Yshi",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "334"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Spirit),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature.Neighboring(),
		Granted: card.FightReap(card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
	}),
)
