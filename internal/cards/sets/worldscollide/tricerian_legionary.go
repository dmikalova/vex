package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tricerian Legionary
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Taunt.
//	Play: Ward a friendly creature.
var TricerianLegionary = card.New(
	"Tricerian Legionary",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "197"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAbility(
		card.Trigger.Play, card.Ward{Target: card.Target.FriendlyCreature}),
)
