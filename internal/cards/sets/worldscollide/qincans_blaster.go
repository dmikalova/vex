package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Qincan's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Sci. Officer Qincan, and you may archive a creature from play."
var QincansBlaster = card.New(
	"Qincan's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Rare,
	card.Provenance(card.WC, "351"),
	card.WithAemberBonus(1),
	card.Connects(card.Pull(SciOfficerQincan, 1)),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "Sci. Officer Qincan"},
				card.May{Do: card.ArchiveFromPlay{Target: card.Target.Creature}},
			}},
		}}}),
	}),
)
