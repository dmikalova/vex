package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Follow the Leader
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For the remainder of the turn, each friendly creature may fight.
var FollowTheLeader = card.New(
	"Follow the Leader",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 8),
	card.WithAbility(card.Trigger.Play, card.GrantFightAnyHouse{}),
)
