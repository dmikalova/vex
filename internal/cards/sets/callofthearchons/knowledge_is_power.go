package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Knowledge is Power
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose one:
//	- Archive a card from your hand
//	- For each card in your archives, gain 1 Æmber.
var KnowledgeIsPower = card.New(
	"Knowledge is Power",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 113),
	card.WithAbility(
		card.Trigger.Play, card.ChooseOne{Options: []card.Effect{
			card.ArchiveFromHand{Count: 1},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.CardsInArchives{Player: card.Controller},
			},
		}}),
)
