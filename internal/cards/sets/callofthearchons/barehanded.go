package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Barehanded
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Put each artifact on top of its owner's deck.
var Barehanded = card.New(
	"Barehanded",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 2),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PutFromPlay{
			Target:      card.Target.EachArtifact,
			Destination: card.To.TopOfDeck,
		}),
)
