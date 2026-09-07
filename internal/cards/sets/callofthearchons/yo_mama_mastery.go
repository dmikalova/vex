package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Yo Mama Mastery
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains taunt.
//	Play: Fully heal this creature.
var YoMamaMastery = card.New(
	"Yo Mama Mastery",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "52"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{Keywords: card.Keywords(card.Keyword.Taunt)}),
	card.WithAbility(
		card.Trigger.Play, card.Heal{
			Fully:  true,
			Target: card.Target.This,
		}),
)
