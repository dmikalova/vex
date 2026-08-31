package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Imperial Traitor
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Reveal your opponent's hand, and you may purge a Sanctum card from your opponent's hand.
var ImperialTraitor = card.New(
	"Imperial Traitor",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 272),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.RevealHand{Player: card.Opponent},
				card.PurgeFromHand{
					Player: card.Opponent,
					House:  card.House.Sanctum,
				},
			},
		}),
)
