package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Opportunist
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains elusive.
//	Play: this creature captures 1 Æmber from your opponent.
var Opportunist = set.New(
	"Opportunist",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.MM, "254"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Keywords: card.Keywords(card.Keyword.Elusive),
	}),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
