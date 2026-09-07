package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ingram's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Medic Ingram, and fully heal a creature."
var IngramsBlaster = card.New(
	"Ingram's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Rare — handle manually
	card.Rarity.Rare,
	card.Provenance(card.WC, 348),
	card.WithAemberBonus(1),
	card.Connects(card.Pull(MedicIngram, 1)),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "Medic Ingram"},
				card.Heal{Fully: true, Target: card.Target.Creature},
			}},
		}}}),
	}),
)
