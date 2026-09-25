package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Horseman of Death
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Connected
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play: Put each Horseman creature from your discard pile into your hand.
var HorsemanOfDeath = set.New(
	"Horseman of Death",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.CotA, "246"),
	card.InCluster(horsemenCluster),
	card.WithPower(5),
	card.WithTraits(card.Traits.Horseman, card.Traits.Spirit),
	card.WithAbility(
		card.Trigger.Play, card.PutCard{Zones: []card.Zone{card.Discard},
			Selection: card.Each{
				Type:  card.Type.Creature,
				Trait: card.Traits.Horseman,
			},
			Destination: card.To.Hand,
		}),
)
