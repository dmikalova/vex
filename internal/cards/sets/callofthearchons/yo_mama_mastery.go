package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Yo Mama Mastery
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains taunt.
//	Play: Fully heal this creature.
var YoMamaMastery = set.New(
	"Yo Mama Mastery",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "52"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{Keywords: card.Keywords(card.Keyword.Taunt)}),
	card.WithAbility(
		card.Trigger.Play, card.Heal{
			Fully:  true,
			Target: card.Target.This,
		}),
)
