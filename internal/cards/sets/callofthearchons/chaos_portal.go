package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Chaos Portal
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Choose a house - reveal the top card of your deck. If it is of the chosen house, play it.
var ChaosPortal = card.New(
	"Chaos Portal",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 127),
	card.WithTraits("Location"),
	card.WithAbility(
		card.Trigger.Action, card.ChooseHouseThen{
			Then: card.Sentences{
				Effects: []card.Effect{
					card.RevealTopOfDeck{},
					card.Conditional{
						Cond: card.ItIsOfHouse{House: card.TheChosenHouse},
						Then: card.PlayRevealedCard{},
					},
				},
			},
		}),
)
