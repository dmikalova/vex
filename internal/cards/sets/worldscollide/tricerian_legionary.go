//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tricerian Legionary
var TricerianLegionary = card.New(
	"Tricerian Legionary",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 197),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAbility(
		card.Trigger.Play, card.Ward{Target: card.Target.FriendlyCreature}),
)
