package massmutation

import "github.com/dmikalova/vex/internal/card"

// Dreadbone Decimus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Dinosaur • Assassin
//
//	Play/Fight: You may exalt Dreadbone Decimus -> destroy a creature with lower power than Dreadbone Decimus.
var DreadboneDecimus = set.New(
	"Dreadbone Decimus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "204"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Assassin),
	card.WithAbility(card.Trigger.PlayFight, card.May{Do: card.Then{
		First: card.Exalt{
			Target: card.Target.This,
			Amount: 1,
		},
		Result: card.Destroy{
			Target: card.Target.Creature.Refine(card.PowerLessThanSource()),
		},
	}}),
)
