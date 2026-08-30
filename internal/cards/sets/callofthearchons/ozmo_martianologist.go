package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

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
var Ozmo = card.New(
	"Ozmo, Martianologist",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 148),
	card.WithPower(2),
	card.WithTraits("Human", "Scientist"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithFightOrReap(card.ChooseOne{
		Options: []card.Effect{
			card.Heal{
				Amount: 3,
				Target: card.Target.Creature.OfHouse(card.House.Mars),
			},
			card.Stun{Target: card.Target.Creature.OfHouse(card.House.Mars)},
		},
	}),
)
