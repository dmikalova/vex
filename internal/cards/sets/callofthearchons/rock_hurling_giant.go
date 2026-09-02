package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Rock-Hurling Giant
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Giant
//
//	After you discard a card from your hand, if it is a Brobnar card, you may deal 4 damage to a creature.
var RockHurlingGiant = card.New(
	"Rock-Hurling Giant",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 44),
	card.WithPower(6),
	card.WithTraits("Giant"),
	card.WithAbility(card.Trigger.AfterDiscardFromHand, card.Conditional{
		Cond: card.ItIs{House: card.House.Self},
		Then: card.May{
			Do: card.DealDamage{
				Target: card.Target.Creature,
				Amount: 4,
			},
		},
	}),
)
