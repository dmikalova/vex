package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Low Dawn
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If there are 3 or more Untamed creatures in your discard pile, gain 2 Æmber. Shuffle each Untamed creature from your discard pile into your deck.
var LowDawn = card.New(
	"Low Dawn",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "377"),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.Sentences{Effects: []card.Effect{
		card.Conditional{
			Cond: card.CardsInDiscardAtLeast{
				House:  card.House.Self,
				Type:   card.Type.Creature,
				Amount: 3,
			},
			Then: card.GainAember{Player: card.Controller, Amount: 2},
		},
		card.ShuffleMatchingFromDiscardIntoDeck{
			House: card.House.Self,
			Type:  card.Type.Creature,
		},
	}}),
)
