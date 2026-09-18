package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Ransack
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Steal 1 Æmber, and discard the top card of your deck -> if the discarded card is a Shadows card, repeat this effect.
var Ransack = set.New(
	"Ransack",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "272"),
	card.WithAbility(card.Trigger.Play, card.Repeat{
		Do: card.Sequence{Effects: []card.Effect{
			card.StealAember{Amount: 1},
			card.DiscardTop{Player: card.Controller},
		}},
		Gate: card.While{Cond: card.ItIs{
			House: card.Houses.Named(card.House.Self),
			Noun:  card.ItNoun.DiscardedCard,
		}},
	}),
)
