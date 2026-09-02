package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Invasion Portal
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Discard cards from the top of your deck until you discard a Mars creature or run out of cards -> put it into your hand.
var InvasionPortal = card.New("Invasion Portal",
	card.House.Mars, card.Type.Artifact, card.Rarity.Rare,
	card.Provenance(card.CotA, 185),
	card.WithTraits("Location"),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.DiscardDeckUntil{
				Type:  card.Type.Creature,
				House: card.House.Mars,
			},
			Result: card.PutDiscardedIntoHand{},
		}),
)
