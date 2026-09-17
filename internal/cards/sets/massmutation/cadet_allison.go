package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Cadet Allison
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  8
//	Traits: Human
//
//	Play/Reap: Discard a random card from your hand -> its house becomes your active house.
var CadetAllison = set.Gigantic(
	"Cadet Allison",
	card.House.StarAlliance,
	card.Rarity.Rare,
	card.Provenance(card.MoMu, "318"),
	card.WithPower(8),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(
		card.Trigger.PlayReap, card.Then{
			First: card.DiscardCard{
				Player:    card.Controller,
				Zones:     []card.Zone{card.Hand},
				Selection: card.Random{Count: 1},
			},
			Result: card.ChangeActiveHouse{To: card.TheContextualHouse},
		}),
)
