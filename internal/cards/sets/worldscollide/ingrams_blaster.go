package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ingram's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Ingram's Blaster to Medic Ingram -> fully heal a creature."
var IngramsBlaster = card.New(
	"Ingram's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "348"),
	card.Connects(card.Pull(MedicIngram, 1)),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: MedicIngram.Name},
				Result: card.Heal{Fully: true, Target: card.Target.Creature},
			},
		}}),
	}),
)
