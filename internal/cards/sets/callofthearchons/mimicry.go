package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mimicry
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Play a Tactic from your opponent's discard pile.
var Mimicry = set.New(
	"Mimicry",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "328"),
	card.WithAbility(
		card.Trigger.Play, card.PlayFrom{
			From:   card.Discard,
			Player: card.Opponent,
			Types:  card.Types.Of(card.Type.Tactic),
		}),
)
