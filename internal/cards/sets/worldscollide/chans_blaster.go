package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Chan's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Chan's Blaster to Commander Chan -> use another creature."
var ChansBlaster = card.New(
	"Chan's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "345"),
	card.Connects(card.Pull(CommanderChan, 1)),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: CommanderChan.Name},
				Result: card.Use{Max: 1, Target: card.Target.OtherFriendlyCreature},
			},
		}}),
	}),
)
