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
var IronObelisk = card.New(
	"Iron Obelisk",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 23),
	card.WithTraits("Location"),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 1).Per(card.InPlay{
		Player:  card.Controller,
		Type:    card.Type.Creature,
		House:   card.House.Brobnar,
		Damaged: true,
	})),
)
