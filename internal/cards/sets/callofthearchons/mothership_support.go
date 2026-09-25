package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Mothership Support
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: For each friendly ready Mars creature, deal 2 damage to a creature.
var MothershipSupport = set.New("Mothership Support",
	card.House.Mars, card.Type.Tactic, card.Rarity.Uncommon,
	card.Provenance(card.CotA, "171"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ForEach{
			Times: card.CardsInPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.Houses.Named(card.House.Self),
				Ready:  true,
			},
			Do: card.DealDamage{
				Target: card.Target.Creature,
				Amount: 2,
			},
		}),
)
