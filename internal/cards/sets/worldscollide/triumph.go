package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Triumph
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: If there are no enemy creatures in play, exalt each friendly creature. If there are 6 or more friendly creatures in play, forge a key at no cost -> purge Triumph.
var Triumph = set.New(
	"Triumph",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "234"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.CardsInPlay{
				Player: card.Opponent,
				Type:   card.Type.Creature,
				None:   true,
			},
			Then: card.Sentences{Effects: []card.Effect{
				card.Exalt{
					Target: card.Target.EachFriendlyCreature,
					Amount: 1,
				},
				card.Conditional{
					Cond: card.CardsInPlay{
						Player: card.Controller,
						Type:   card.Type.Creature,
						Amount: 6,
					},
					Then: card.ForgeKey{FreeOfCost: true},
				},
			}},
		}),
)
