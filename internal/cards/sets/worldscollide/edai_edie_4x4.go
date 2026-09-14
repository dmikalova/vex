package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// EDAI "Edie" 4x4
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: AI • Scientist
//
//	Your opponent's keys cost +1 Æmber for each card in your archives.
//	Play: Archive a card from your hand.
var EDAIEdie4x4 = set.New(
	"EDAI \"Edie\" 4x4",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "132"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Ai, card.Traits.Scientist),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 1).Per(card.CardsInArchives{
		Player: card.Controller,
	})),
	card.WithAbility(
		card.Trigger.Play,
		card.ArchiveCard{Zone: card.Hand, Selection: card.Chosen{}},
	),
)
