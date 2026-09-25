package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Tribune Pompitus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  2
//	Traits: Dinosaur • Politician
//
//	Each friendly creature gains +2 power for each Æmber on it.
//	Before Fight: You may exalt Tribune Pompitus.
var TribunePompitus = set.New(
	"Tribune Pompitus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "213"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithConstant(card.ConstantAbility{
		Target:     card.Target.EachFriendlyCreature,
		PowerBonus: 2,
		PerTarget:  card.AemberOnIt,
	}),
	card.WithAbility(
		card.Trigger.BeforeFight, card.May{Do: card.Exalt{
			Target: card.Target.This,
			Amount: 1,
		}}),
)
