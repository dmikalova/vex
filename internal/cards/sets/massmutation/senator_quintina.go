package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Senator Quintina
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Dinosaur • Politician
//
//	After a creature reaps, exalt it.
var SenatorQuintina = set.New(
	"Senator Quintina",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "211"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Exalt{
			Target: card.Target.Triggering,
			Amount: 1,
		}),
)
