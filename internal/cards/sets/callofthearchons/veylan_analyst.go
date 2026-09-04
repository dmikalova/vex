package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Veylan Analyst
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	After you use a card, if it is an artifact, gain 1 Æmber.
var VeylanAnalyst = card.New(
	"Veylan Analyst",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 156),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithAbility(card.Trigger.AfterUse, card.Conditional{
		Cond: card.ItIs{Type: card.Type.Artifact},
		Then: card.GainAember{
			Player: card.Controller,
			Amount: 1,
		},
	}),
)
