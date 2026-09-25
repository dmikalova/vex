package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Imperial Traitor
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Reveal your opponent's hand. You may purge a Sanctum card from your opponent's hand.
var ImperialTraitor = set.New(
	"Imperial Traitor",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "272"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.RevealHand{Player: card.Opponent},
				card.PurgeCard{
					Zones:  []card.Zone{card.Hand},
					Player: card.Opponent,
					Selection: card.Chosen{
						House:    card.Houses.Named(card.House.Sanctum),
						Optional: true,
					},
				},
			},
		}),
)
