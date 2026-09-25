package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Ozmo, Martianologist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Scientist
//
//	Elusive.
//	Fight/Reap: Choose one:
//	- Heal 3 damage from a Mars creature
//	- Stun a Mars creature.
var Ozmo = set.New(
	"Ozmo, Martianologist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "148"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(card.Trigger.FightReap, card.ChooseOne{
		Options: []card.Effect{
			card.Heal{
				Amount: 3,
				Target: card.Target.Creature.House(card.Houses.Named(card.House.Mars)),
			},
			card.Stun{Target: card.Target.Creature.House(card.Houses.Named(card.House.Mars))},
		},
	}),
)
