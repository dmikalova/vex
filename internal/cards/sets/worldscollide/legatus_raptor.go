package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Legatus Raptor
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Fight: You may exalt Legatus Raptor. Ready and use another friendly creature.
var LegatusRaptor = set.New(
	"Legatus Raptor",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "187"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.Fight, card.May{Do: card.Sequence{Effects: []card.Effect{
			card.Exalt{
				Target: card.Target.This,
				Amount: 1,
			},
			card.OnChooseCreature{
				Target: card.Target.OtherFriendlyCreature,
				Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.UseVerb{}},
			},
		}}}),
)
