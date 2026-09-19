package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Sci. Officer Qincan
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Alien • Proximan • Scientist
//
//	Elusive.
//	After a player chooses an active house which matches no cards in play, steal 1 Æmber.
var SciOfficerQincan = set.New(
	"Sci. Officer Qincan",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "304"),
	card.InCluster(card.Pulled(qincansBlasterCluster, 1, 1.25)),
	card.WithPower(2),
	card.WithTraits(card.Traits.Alien, card.Traits.Proximan, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithEachPlayerAbility(
		card.Trigger.AfterChooseHouse, card.Conditional{
			Cond: card.ActiveHouseMatchesNoCardsInPlay{},
			Then: card.StealAember{Amount: 1},
		}),
)
