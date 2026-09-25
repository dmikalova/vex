package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Stunner
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: You may stun a creature."
var Stunner = set.New(
	"Stunner",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "319"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.May{Do: card.Stun{Target: card.Target.Creature}}),
	}),
)
