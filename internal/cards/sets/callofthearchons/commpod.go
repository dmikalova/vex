package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Commpod
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Reveal any number of Mars cards from your hand, and for each card revealed this way, ready a friendly Mars Creature.
var Commpod = set.New(
	"Commpod",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "181"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{
			Effects: []card.Effect{
				card.RevealHand{
					Player: card.Controller,
					House:  card.Houses.Named(card.House.Self),
				},
				card.ReadyCreatures{
					Max: card.CardsRevealed{},
					Target: card.Target.EachFriendlyCreature.House(
						card.Houses.Named(card.House.Self),
					),
				},
			},
		}),
)
