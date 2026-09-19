package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Paraguardian
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Reap: You may exalt and ward Paraguardian.
var Paraguardian = set.New(
	"Paraguardian",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "206"),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.Reap, card.May{Do: card.Sequence{Effects: []card.Effect{
			card.Exalt{
				Target: card.Target.This,
				Amount: 1,
			},
			card.Ward{Target: card.Target.This.NeighborsOf()},
		}}}),
)
