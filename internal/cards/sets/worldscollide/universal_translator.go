package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Universal Translator
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Use a non-Star Alliance creature."
var UniversalTranslator = card.New(
	"Universal Translator",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "322"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.Use{
			Max:    1,
			Target: card.Target.EachFriendlyCreature.ExceptHouse(card.House.Self),
		}),
	}),
)
