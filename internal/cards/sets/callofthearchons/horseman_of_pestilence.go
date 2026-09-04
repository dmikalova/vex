package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Horseman of Pestilence
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play/Fight/Reap: Deal 1 damage to each non-Horseman trait creature.
var HorsemanOfPestilence = card.New(
	"Horseman of Pestilence",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 248),
	card.WithPower(5),
	card.WithTraits(card.Traits.Horseman, card.Traits.Spirit),
	card.WithPlayFightReap(card.DealDamage{
		Amount: 1,
		Target: card.Target.EachCreature.ExceptTrait(card.Traits.Horseman),
	}),
)
