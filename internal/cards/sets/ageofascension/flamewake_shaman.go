package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Flamewake Shaman
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Human
//
//	Play: Deal 2 damage to a creature.
var FlamewakeShaman = card.New(
	"Flamewake Shaman",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 22),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.Creature,
		}),
)
