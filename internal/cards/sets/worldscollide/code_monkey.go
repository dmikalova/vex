package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Code Monkey
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: AI • Beast
//
//	Deploy.
//	Play: Archive each neighboring creature from play. If those creatures share a house, gain 2 Æmber.
var CodeMonkey = set.New(
	"Code Monkey",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "147"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Ai, card.Traits.Beast),
	card.WithKeywords(card.Keyword.Deploy),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.ArchiveFromPlay{
				Target: card.Target.EachCreature.Neighboring(),
			},
			card.Conditional{
				Cond: card.ArchivedCreaturesShareHouse{},
				Then: card.GainAember{Player: card.Controller, Amount: 2},
			},
		}}),
)
