package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Draco Praeco
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Dinosaur • Politician
//
//	Reap: You may exalt Draco Praeco, and choose a house - enrage each creature of the chosen house.
var DracoPraeco = card.New(
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
				Then: card.Enrage{Target: card.Target.EachCreature.OfChosenHouse()},
			},
		}}}),
)
