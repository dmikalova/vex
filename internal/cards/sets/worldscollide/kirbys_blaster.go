package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Kirby's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Com. Officer Kirby, and draw 2 cards."
var KirbysBlaster = card.New(
	"Kirby's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "350"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "Com. Officer Kirby"},
				card.Draw{Amount: 2},
			}},
		}}}),
	}),
)
