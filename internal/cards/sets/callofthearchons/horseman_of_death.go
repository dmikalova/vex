package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Horseman of Death
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: FIXED
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play: Put each Horseman trait creature from your discard pile into your hand.
var HorsemanOfDeath = card.New(
	"Horseman of Death",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Fixed,
	card.Provenance(card.CotA, 246),
	card.WithPower(5),
	card.WithTraits("Horseman", "Spirit"),
	card.WithAbility(
		card.Trigger.Play, card.PutFromDiscard{
			Type:        card.Type.Creature,
			Trait:       "Horseman",
			All:         true,
			Destination: card.To.Hand,
		}),
)
