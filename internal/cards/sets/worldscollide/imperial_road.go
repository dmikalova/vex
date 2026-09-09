package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Imperial Road
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Versatile.
//	Action: Play a Saurian creature -> stun it.
var ImperialRoad = card.New(
	"Imperial Road",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "223"),
	card.WithTraits(card.Traits.Location),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.PlayFrom{
				From:  card.Hand,
				House: card.House.Self,
				Type:  card.Type.Creature,
			},
			Result: card.Stun{Target: card.Target.Triggering},
		}),
)
