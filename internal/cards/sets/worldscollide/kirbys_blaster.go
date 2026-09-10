package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Kirby's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Kirby's Blaster to Com. Officer Kirby -> draw 2 cards."
var KirbysBlaster = card.New(
	"Kirby's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "350"),
	card.Connects(card.Pull(ComOfficerKirby, 1)),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: ComOfficerKirby.Name},
				Result: card.Draw{Amount: 2},
			},
		}}),
	}),
)
