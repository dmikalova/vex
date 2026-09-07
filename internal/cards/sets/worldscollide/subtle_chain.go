package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Subtle Chain
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Your opponent discards a random card from their hand.
var SubtleChain = card.New(
	"Subtle Chain",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 262),
	// TODO(duplicate): mechanically identical to Mind Barb (Dis) — fold/handle manually.
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DiscardRandomFromHand{Player: card.Opponent}),
)
