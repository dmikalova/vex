package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Chaos Portal
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Choose a house. Reveal the top card of your deck. If it is of the chosen house, play it.
var ChaosPortal = set.New(
	"Chaos Portal",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "127"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.ChooseHouseThen{
			Then: card.Sequence{
				Effects: []card.Effect{
					card.RevealTopOfDeck{Amount: 1},
					card.Conditional{
						Cond: card.ItIs{House: card.Houses.Chosen},
						Then: card.PlayRevealedCard{},
					},
				},
			},
		}),
)
