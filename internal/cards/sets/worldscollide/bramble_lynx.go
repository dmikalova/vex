package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Bramble Lynx
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Skirmish.
//	If you have used a Creature to reap this turn, Bramble Lynx enters play ready.
var BrambleLynx = card.New(
	"Bramble Lynx",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "353"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithEntersPlay(card.Conditional{
		Cond: card.UsedCreatureToReap{},
		Then: card.Ready{Target: card.Target.This},
	}),
)
