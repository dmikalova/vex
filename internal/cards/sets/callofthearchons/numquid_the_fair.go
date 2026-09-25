package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Numquid the Fair
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	Play: Destroy an enemy creature -> if you are overwhelmed, repeat this effect.
var NumquidTheFair = set.New(
	"Numquid the Fair",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "253"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(card.Trigger.Play, card.Repeat{
		Do:   card.Destroy{Target: card.Target.EnemyCreature},
		Gate: card.While{Cond: card.Overwhelmed{}},
	}),
)
