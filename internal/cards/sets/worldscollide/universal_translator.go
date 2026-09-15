package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Universal Translator
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Use a non-Star Alliance creature."
var UniversalTranslator = set.New(
	"Universal Translator",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "322"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.Use{
			Max:    1,
			Target: card.Target.EachFriendlyCreature.House(card.Houses.Except(card.House.Self)),
		}),
	}),
)
