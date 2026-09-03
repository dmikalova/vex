package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Unguarded Camp
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: For each creature you have in excess of your opponent, a friendly creature captures 1 Æmber from your opponent. Each creature cannot capture more than 1 Æmber this way.
var UnguardedCamp = card.New(
	"Unguarded Camp",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 17),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount:   1,
			Target:   card.Target.FriendlyCreature,
			Source:   card.Opponent,
			Per:      card.ExcessCreatures{Player: card.Controller},
			Distinct: true,
		}),
)
