package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Velum
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Connected
//	Power:  2
//	Traits: Human • Scientist
//
//	Reap: Archive a card from your hand. If you control Hyde, archive a card from your hand.
//	Destroyed: Archive Hyde from your discard pile -> archive Velum from play.
var Velum = card.New(
	"Velum",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "181"),
	card.InCluster(hydeCluster),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(card.Trigger.Reap, card.Sentences{Effects: []card.Effect{
		card.ArchiveCard{Zone: card.Hand, Selection: card.Chosen{}},
		card.Conditional{
			Cond: card.ControlsNamed{Name: HydeName},
			Then: card.ArchiveCard{Zone: card.Hand, Selection: card.Chosen{}},
		},
	}}),
	card.WithAbility(card.Trigger.Destroyed, card.Then{
		First:  card.ArchiveCard{Zone: card.Discard, Selection: card.Named{Name: HydeName}},
		Result: card.ArchiveFromPlay{Target: card.Target.This},
	}),
)
