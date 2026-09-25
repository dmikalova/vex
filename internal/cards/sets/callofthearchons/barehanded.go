package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Barehanded
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Put each artifact on top of its owner's deck.
var Barehanded = set.New(
	"Barehanded",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "2"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.PutFromPlay{
			Target:      card.Target.EachArtifact,
			Destination: card.To.TopOfDeck,
		}),
)
