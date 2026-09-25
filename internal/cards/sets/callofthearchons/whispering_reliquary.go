package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Whispering Reliquary
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Put an artifact into its owner's hand.
var WhisperingReliquary = set.New(
	"Whispering Reliquary",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "237"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.PutFromPlay{
			Target:      card.Target.Artifact,
			Destination: card.To.Hand,
		}),
)
