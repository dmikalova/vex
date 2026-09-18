package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Iron Obelisk
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Your opponent's keys cost +1 Æmber for each friendly damaged Brobnar creature.
var IronObelisk = set.New(
	"Iron Obelisk",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "23"),
	card.WithTraits(card.Traits.Location),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 1).Per(card.CardsInPlay{
		Player:  card.Controller,
		Type:    card.Type.Creature,
		House:   card.Houses.Named(card.House.Self),
		Damaged: true,
	})),
)
