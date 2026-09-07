package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Khrkhar's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Lieutenant Khrkhar, and ward Lieutenant Khrkhar."
var KhrkharsBlaster = card.New(
	"Khrkhar's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "349"),
	card.WithAemberBonus(1),
	card.Connects(card.Pull(LieutenantKhrkhar, 1)),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "Lieutenant Khrkhar"},
				card.Ward{Target: card.Target.AttachedHost.Named("Lieutenant Khrkhar")},
			}},
		}}}),
	}),
)
