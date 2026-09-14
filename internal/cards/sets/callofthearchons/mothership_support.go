package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mothership Support
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each friendly ready Mars Creature, deal 2 damage to a Creature.
var MothershipSupport = set.New("Mothership Support",
	card.House.Mars, card.Type.Tactic, card.Rarity.Uncommon,
	card.Provenance(card.CotA, "171"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ForEach{
			Times: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.House.Self,
				Ready:  true,
			},
			Do: card.DealDamage{
				Target: card.Target.Creature,
				Amount: 2,
			},
		}),
)
