package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Invasion Portal
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Discard cards from the top of your deck until you discard a Mars Creature or run out of cards -> put the discarded Creature into your hand.
var InvasionPortal = set.New("Invasion Portal",
	card.House.Mars, card.Type.Artifact, card.Rarity.Rare,
	card.Provenance(card.CotA, "185"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.DiscardTopOfDeckUntil{
				Type:  card.Type.Creature,
				House: card.House.Self,
			},
			Result: card.PutDiscardedIntoHand{Type: card.Type.Creature},
		}),
)
