package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Molina's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Armsmaster Molina, and you may deal 3 damage to a creature."
var MolinasBlaster = card.New(
	"Molina's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Rare — handle manually
	card.Rarity.Rare,
	card.Provenance(card.WC, 302),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "Armsmaster Molina"},
				card.May{Do: card.DealDamage{Amount: 3, Target: card.Target.Creature}},
			}},
		}}}),
	}),
)
