//go:build todo

// TODO: blocked on two gaps before this can build:
//  1. source rarity is "Variant"; awaiting a real-rarity mapping (as the Brews).
//  2. the payoff "stun a creature for each upgrade on Chief Engineer Walls" needs
//     a Count over the upgrades attached to a creature, which no effect exposes
//     yet. The Fight/Reap choice and the attach are wired below; the stun payoff
//     is left as a TODO until that Count lands.
package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Walls' Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may deal 2 damage to a creature, or attach Walls' Blaster to Chief Engineer Walls."
//	After you attach Walls' Blaster to Chief Engineer Walls, stun a creature for each upgrade on Chief Engineer Walls.
var WallsBlaster = card.New(
	"Walls' Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "352"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "Chief Engineer Walls"},
				// TODO: stun a creature for each upgrade on Chief Engineer Walls,
				// once a Count over a creature's upgrades exists.
			}},
		}}}),
	}),
)
