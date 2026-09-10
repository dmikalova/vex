package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Molina's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Molina's Blaster to Armsmaster Molina -> deal 3 damage to a creature."
var MolinasBlaster = card.New(
	"Molina's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "302"),
	card.Connects(card.Pull(ArmsmasterMolina, 1)),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: ArmsmasterMolina.Name},
				Result: card.DealDamage{Amount: 3, Target: card.Target.Creature},
			},
		}}),
	}),
)
