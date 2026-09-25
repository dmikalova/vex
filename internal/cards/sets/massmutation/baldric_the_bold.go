package massmutation

import "github.com/dmikalova/vex/internal/card"

// Baldric the Bold
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Before Fight: If the fights creature is the most powerful enemy creature, gain 2 Æmber.
var BaldricTheBold = set.New(
	"Baldric the Bold",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "144"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithAbility(card.Trigger.BeforeFight, card.Conditional{
		Cond: card.ItIsAmong{
			Target: card.Target.EachEnemyCreature.Refine(card.MostPowerful),
			Noun:   card.ItNoun.FoughtCreature,
		},
		Then: card.GainAember{
			Player: card.Controller,
			Amount: 2,
		},
	}),
)
