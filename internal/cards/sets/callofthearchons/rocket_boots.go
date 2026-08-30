package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Rocket Boots
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains, "Fight/Reap: If this is the first time this creature was used this turn, ready it."
var RocketBoots = card.New(
	"Rocket Boots",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 158),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.ReadyIfFirstUse{Target: card.Target.This}),
	}),
)
