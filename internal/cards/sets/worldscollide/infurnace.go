package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Infurnace
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Play: Purge up to 2 cards from a discard pile. Your opponent loses Æmber equal to the total Æmber bonus of the purged cards.
var Infurnace = set.New(
	"Infurnace",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "78"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.PurgeCard{
				Player:    card.ChosenPlayer,
				Selection: card.Chosen{Optional: true},
				Amount:    2,
			},
			card.LoseAember{
				Player:  card.Opponent,
				EqualTo: card.PurgedAemberBonus{},
			},
		}},
	),
)
