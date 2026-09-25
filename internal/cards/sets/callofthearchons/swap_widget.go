package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Swap Widget
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Put a friendly ready Mars creature into its owner's hand -> put a Mars creature with a different name from your hand into play. Ready it.
var SwapWidget = set.New(
	"Swap Widget",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "189"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.PutFromPlay{
				Target: card.Target.FriendlyCreature.House(card.Houses.Named(card.House.Self)).
					Ready(),
				Destination: card.To.Hand,
			},
			Result: card.Sequence{Effects: []card.Effect{
				card.PutFromHand{
					Type:           card.Type.Creature,
					House:          card.Houses.Named(card.House.Self),
					ExceptSameName: true,
				},
				card.Ready{Target: card.Target.Triggering},
			}},
		}),
)
