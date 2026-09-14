package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Titan Guardian
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Armor:  1
//	Traits: Beast • Cyborg
//
//	Taunt.
//	Destroyed: If Titan Guardian is not on a flank, draw 2 cards.
var TitanGuardian = set.New(
	"Titan Guardian",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "141"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Beast, card.Traits.Cyborg),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAbility(
		card.Trigger.Destroyed, card.Conditional{
			Cond: card.Not{Cond: card.OnFlank{}},
			Then: card.Draw{Amount: 2},
		}),
)
