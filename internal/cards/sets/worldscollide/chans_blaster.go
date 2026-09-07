package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Chan's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Commander Chan, and you may use an another creature."
var ChansBlaster = card.New(
	"Chan's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "345"),
	card.WithAemberBonus(1),
	card.Connects(card.Pull(CommanderChan, 1)),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "Commander Chan"},
				card.May{Do: card.Use{Max: 1, Target: card.Target.OtherFriendlyCreature}},
			}},
		}}}),
	}),
)
