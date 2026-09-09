package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hyde
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Human • Scientist
//
//	Reap: Draw a card. If you control Velum, draw a card.
//	Destroyed: Archive Velum from your discard pile -> archive Hyde from play.
var Hyde = card.New(
	"Hyde",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "167"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.Connects(card.PullExact(Velum, 1)),
	card.WithAbility(card.Trigger.Reap, card.Sentences{Effects: []card.Effect{
		card.Draw{Amount: 1},
		card.Conditional{
			Cond: card.ControlsNamed{Name: "Velum"},
			Then: card.Draw{Amount: 1},
		},
	}}),
	card.WithAbility(card.Trigger.Destroyed, card.Then{
		First:  card.ArchiveFromDiscard{Name: "Velum"},
		Result: card.ArchiveFromPlay{Target: card.Target.This},
	}),
)
