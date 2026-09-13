package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Doctor Driscoll
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Scientist
//
//	Elusive.
//	Action: Heal 2 damage from a Creature. For each damage healed this way, gain 1 Æmber.
var DoctorDriscoll = card.New(
	"Doctor Driscoll",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "329"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Heal{
				Amount: 2,
				Target: card.Target.Creature,
			},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.DamageHealed{},
			},
		}}),
)
