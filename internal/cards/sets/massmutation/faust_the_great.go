package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Faust the Great
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Dinosaur
//
//	Your opponent's keys cost +1 Æmber for each friendly creature with Æmber on it.
//	Play: You may exalt a friendly creature.
var FaustTheGreat = set.New(
	"Faust the Great",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "192"),
	card.InCluster(card.Pulled(monumentToFaustCluster, 1, 1)),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 1).Per(card.CardsInPlay{
		Player:     card.Controller,
		Type:       card.Type.Creature,
		WithAember: true,
	})),
	card.WithAbility(
		card.Trigger.Play, card.May{
			Do: card.Exalt{Target: card.Target.FriendlyCreature, Amount: 1},
		}),
)
