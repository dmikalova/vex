package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Eyegor
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg
//
//	Play: Look at the top 3 cards of your deck, put 1 into your hand, and discard 2.
var Eyegor = set.New(
	"Eyegor",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "111"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg),
	card.WithAbility(
		card.Trigger.Play, card.LookAtTopOfDeck{
			Amount: 3,
			Then: []card.TopAct{
				card.ChooseAndMove{
					Cards: 1,
					Dest:  card.Into.Hand,
				},
				card.ChooseAndMove{
					Cards: 2,
					Dest:  card.Into.Discard,
				},
			},
		}),
)
