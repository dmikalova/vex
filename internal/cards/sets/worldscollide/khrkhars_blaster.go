package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Khrkhar's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Khrkhar's Blaster to Lieutenant Khrkhar -> ward Lieutenant Khrkhar."
var KhrkharsBlaster = card.New(
	"Khrkhar's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "349"),
	card.Connects(card.Pull(LieutenantKhrkhar, 1)),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: LieutenantKhrkhar.Name},
				Result: card.Ward{Target: card.Target.AttachedHost.Named(LieutenantKhrkhar.Name)},
			},
		}}),
	}),
)
