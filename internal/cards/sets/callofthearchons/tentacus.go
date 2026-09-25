package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Tentacus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	In order to use an artifact, your opponent must give you 1 Æmber.
var Tentacus = set.New(
	"Tentacus",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "100"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithRestrictions(card.Restrictions{
		Toll: card.Toll{
			Action: card.TollOn.UseArtifact,
			Amount: 1,
		},
	}),
)
