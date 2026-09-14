package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Rocket Boots
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This Creature gains, "Fight/Reap: If this is the first time this Creature was used this turn, ready it."
var RocketBoots = set.New(
	"Rocket Boots",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "158"),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ReadyIfFirstUse{Target: card.Target.This}),
	}),
)
