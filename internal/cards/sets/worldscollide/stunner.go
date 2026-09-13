package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Stunner
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains, "Fight/Reap: You may stun a Creature."
var Stunner = card.New(
	"Stunner",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "319"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.May{Do: card.Stun{Target: card.Target.Creature}}),
	}),
)
