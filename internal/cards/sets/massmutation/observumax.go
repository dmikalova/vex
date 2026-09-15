package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Observe-u-Max
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: This creature captures 1 Æmber from your opponent."
var ObservuMax = set.New(
	"Observe-u-Max",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.MM, "309"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
	}),
)
