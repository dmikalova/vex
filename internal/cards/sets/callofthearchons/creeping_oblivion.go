package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Creeping Oblivion
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Purge up to 2 cards from a discard pile.
var CreepingOblivion = card.New(
	"Creeping Oblivion",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 56),
	card.WithAbility(
		card.Trigger.Play, card.PurgeCard{
			Zone:   card.Discard,
			Amount: 2,
			UpTo:   true,
		}),
)
