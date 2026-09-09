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
//	Reap: Draw a card. If you control Velum, draw 2 cards instead.
//	Destroyed: Archive Velum from your discard pile. If you do, archive Hyde.
var Hyde = card.New(
	"Hyde",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "167"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.Connects(card.Pull(Velum, 1)),
	card.WithAbility(card.Trigger.Reap, card.Draw{
		Amount: 1,
		Or: card.OrAmount{
			Amount: 2,
			When:   card.ControlsNamed{Name: "Velum"},
		},
	}),
	card.WithAbility(card.Trigger.Destroyed, card.Then{
		First:  card.ArchiveFromDiscard{Name: "Velum"},
		Result: card.ArchiveFromPlay{Target: card.Target.This},
	}),
)
