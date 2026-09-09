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
var Infurnace = card.New(
	"Infurnace",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "78"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.PurgeCard{Zone: card.Discard, Amount: 2, UpTo: true},
			card.LoseAemberEqualTo{
				Player: card.Opponent,
				Count:  card.PurgedAemberBonus{},
			},
		}},
	),
)
