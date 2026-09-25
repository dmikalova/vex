package massmutation

import "github.com/dmikalova/vex/internal/card"

// Nerotaurus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Dinosaur • Politician
//
//	Fight: Your opponent cannot use creatures to reap during their next turn.
//	Reap: Your opponent cannot use creatures to fight during their next turn.
var Nerotaurus = set.New(
	"Nerotaurus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "225"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithAbility(
		card.Trigger.Fight, card.Restrict{
			Player:   card.Opponent,
			Action:   card.Restricted.Reaping,
			Duration: card.Duration.OpponentNextTurn,
		}),
	card.WithAbility(
		card.Trigger.Reap, card.Restrict{
			Player:   card.Opponent,
			Action:   card.Restricted.Fighting,
			Duration: card.Duration.OpponentNextTurn,
		}),
)
