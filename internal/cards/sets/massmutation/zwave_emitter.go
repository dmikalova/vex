package massmutation

import "github.com/dmikalova/vex/internal/card"

// Z-Wave Emitter
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Connected
//	Bonus:  Æmber
//
//	This creature gains, "At the start of your turn, ward this creature."
var ZWaveEmitter = set.New(
	"Z-Wave Emitter",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Connected,
	card.Provenance(card.MM, "356"),
	card.InCluster(zForceCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.StartOfTurn,
			Effect:  card.Ward{Target: card.Target.This},
		}},
	}),
)
