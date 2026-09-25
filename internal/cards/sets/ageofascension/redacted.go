package ageofascension

import "github.com/dmikalova/vex/internal/card"

// [REDACTED]
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: [redacted]
//
//	After you choose Logos as your active house, place 1 Æmber from the common supply on [REDACTED]. If there are 4 or more Æmber on it, forge a key at no cost -> purge [REDACTED].
var REDACTED = set.New(
	"[REDACTED]",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "139"),
	card.WithTraits(card.Traits.Redacted),
	card.WithAbility(
		card.Trigger.AfterChooseHouse, card.Conditional{
			Cond: card.ChoseHouse{House: card.House.Self},
			Then: card.Sequence{Effects: []card.Effect{
				card.PlaceAemberOnThis{Amount: 1},
				card.Conditional{
					Cond: card.CountIs{
						Count:  card.AemberOnThis{},
						Is:     card.AtLeast,
						Amount: 4,
					},
					Then: card.ForgeKey{FreeOfCost: true},
				},
			}},
		}),
)
