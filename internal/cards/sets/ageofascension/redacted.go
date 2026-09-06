package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// [REDACTED]
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: [redacted]
//
//	After you choose Logos as your active house, place 1 Æmber from the common supply on [REDACTED]. If there are 4 or more Æmber on it, destroy [REDACTED], and forge a key at no cost.
var REDACTED = card.New(
	"[REDACTED]",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 139),
	card.WithTraits(card.Traits.Redacted),
	card.WithAbility(
		card.Trigger.AfterChooseHouse, card.Conditional{
			Cond: card.ChoseHouse{House: card.House.Self},
			Then: card.Sentences{Effects: []card.Effect{
				card.PlaceAemberOnThis{Amount: 1},
				card.Conditional{
					Cond: card.AemberOnThisAtLeast{Amount: 4},
					Then: card.Sequence{Effects: []card.Effect{
						card.Destroy{Target: card.Target.This},
						card.ForgeKey{FreeOfCost: true},
					}},
				},
			}},
		}),
)
