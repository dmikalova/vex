package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bonkers Killing Machine
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Weapon
//
//	Action: Discard the top card of each player's deck. For each card discarded this way, destroy a creature or artifact of that card's house. If fewer than 2 cards are destroyed this way, destroy Bonkers Killing Machine.
var BonkersKillingMachine = set.New(
	"Bonkers Killing Machine",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "128"),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{
			Effects: []card.Effect{
				card.DiscardTop{
					Amount: 1,
					Player: card.EachPlayer,
				},
				card.ForEachDiscarded{
					Do: card.Destroy{
						Target: card.Target.CreatureOrArtifact.House(card.Houses.Contextual),
					},
				},
				card.Conditional{
					Cond: card.Not{Cond: card.CountIs{
						Count:  card.CardsDestroyed{},
						Is:     card.AtLeast,
						Amount: 2,
					}},
					Then: card.Destroy{Target: card.Target.This},
				},
			},
		}),
)
