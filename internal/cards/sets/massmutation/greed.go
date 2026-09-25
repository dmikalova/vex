package massmutation

import "github.com/dmikalova/vex/internal/card"

// Greed
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Demon • Sin
//
//	For each friendly Sin creature your hand size is 1 more.
var Greed = set.New(
	"Greed",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "058"),
	card.InCluster(sinsCluster),
	card.OneCopyPerDeck(),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	card.WithDrawModifierPer(card.Controller, 1, card.CardsInPlay{
		Player: card.Controller,
		Type:   card.Type.Creature,
		Trait:  card.Traits.Sin,
	}),
)
