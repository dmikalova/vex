package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Rocket Boots
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "Fight/Reap: If this is the first time this creature has been used this turn, ready this creature."
var RocketBoots = set.New(
	"Rocket Boots",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "158"),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.Conditional{
			Cond: card.SourceFirstUseThisTurn{},
			Then: card.Ready{Target: card.Target.This},
		}),
	}),
)
