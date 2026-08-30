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
var BonkersKillingMachine = card.New(
	"Bonkers Killing Machine",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 128),
	card.WithTraits("Weapon"),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{
			Effects: []card.Effect{
				card.Sentence{Effect: card.DiscardTopOfEachDeck{}},
				card.Sentence{Effect: card.ForEachDiscarded{
					Do: card.Destroy{Target: card.Target.ChosenInPlay.OfContextualHouse()},
				}},
				card.Conditional{
					Cond: card.CardsDestroyedFewerThan{Amount: 2},
					Then: card.Destroy{Target: card.Target.This},
				},
			},
		}),
)
