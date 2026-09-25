package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Draco Praeco
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Dinosaur • Politician
//
//	Reap: You may exalt Draco Praeco. Choose a house. Enrage each creature of the chosen house.
var DracoPraeco = set.New(
	"Draco Praeco",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "201"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithAbility(
		card.Trigger.Reap, card.May{Do: card.Sequence{Effects: []card.Effect{
			card.Exalt{
				Target: card.Target.This,
				Amount: 1,
			},
			card.ChooseHouseThen{
				Then: card.Enrage{Target: card.Target.EachCreature.House(card.Houses.Chosen)},
			},
		}}}),
)
