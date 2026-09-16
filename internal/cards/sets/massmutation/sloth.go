package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Sloth
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Demon • Sin
//
//	At the end of your turn, if you did not use any creatures this turn, for each friendly Sin creature, gain 1 Æmber.
var Sloth = set.New(
	"Sloth",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "062"),
	card.InCluster(sinsCluster),
	card.OneCopyPerDeck(),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	card.WithAbility(card.Trigger.EndOfTurn, card.Conditional{
		Cond: card.UsedNoCreatures{},
		Then: card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				Trait:  card.Traits.Sin,
			},
		},
	}),
)
