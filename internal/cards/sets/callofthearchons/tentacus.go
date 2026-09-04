package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tentacus
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	Your opponent must give you 1 Æmber in order to use an artifact.
var Tentacus = card.New(
	"Tentacus",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 100),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithRestrictions(card.Restrictions{
		Toll: card.Toll{
			Action: card.TollOn.UseArtifact,
			Amount: 1,
		},
	}),
)
