package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Uxlyx the Zookeeper
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Martian • Scientist
//
//	Elusive.
//	Reap: Put an enemy creature into your archives.
var UxlyxTheZookeeper = card.New(
	"Uxlyx the Zookeeper",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 201),
	card.WithPower(2),
	card.WithTraits("Martian", "Scientist"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.PutFromPlay{
			Target:      card.Target.EnemyCreature,
			Destination: card.To.Archives.Yours(),
		}),
)
