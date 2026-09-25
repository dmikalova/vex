package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Lollop the Titanic
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  11
//	Traits: Giant • Location
//
//	Lollop the Titanic deals no damage when attacked.
var LollopTheTitanic = set.New(
	"Lollop the Titanic",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "14"),
	card.WithPower(11),
	card.WithTraits(card.Traits.Giant, card.Traits.Location),
	card.WithNoDamageWhenAttacked(),
)
