package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Dysania
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant
//
//	Play: For each card in your opponent's archives, gain 1 Æmber. Your opponent discards each of their archived cards.
var Dysania = set.New(
	"Dysania",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "141"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per: card.CardsInZone{
					Zone:   card.Archives,
					Player: card.Opponent,
				},
			},
			card.DiscardArchives{Player: card.Opponent},
		}}),
)
