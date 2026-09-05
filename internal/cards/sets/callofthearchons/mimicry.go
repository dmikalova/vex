package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mimicry
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Play a tactic from your opponent's discard pile.
var Mimicry = card.New(
	"Mimicry",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 328),
	card.WithAbility(
		card.Trigger.Play, card.PlayFrom{
			From:   card.Discard,
			Player: card.Opponent,
			Type:   card.Type.Tactic,
		}),
)
