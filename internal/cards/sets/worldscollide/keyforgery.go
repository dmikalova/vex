package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Keyforgery
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	When your opponent would forge a key, they name a house. Reveal a random card from your hand. If that card is not of the named house, destroy Keyforgery, and they do not forge that key.
var Keyforgery = set.New(
	"Keyforgery",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "271"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.BeforeOpponentForgesKey,
		card.Sentences{Effects: []card.Effect{
			card.OpponentNamesHouse{},
			card.RevealRandomFromHand{},
			card.Conditional{
				Cond: card.ItIsNotOfNamedHouse{Subject: card.Subject.ThatCard},
				Then: card.Sequence{Effects: []card.Effect{
					card.Destroy{Target: card.Target.This},
					card.CancelForge{},
				}},
			},
		}},
	),
)
