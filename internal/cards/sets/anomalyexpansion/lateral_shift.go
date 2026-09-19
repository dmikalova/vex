package anomalyexpansion

import "github.com/dmikalova/vactrol/internal/card"

// Lateral Shift
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//
//	Play: Reveal your opponent's hand. Play a card from your opponent's hand.
var LateralShift = set.New(
	"Lateral Shift",
	card.House.Brobnar,
	card.Type.Tactic,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A03"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.RevealHand{Player: card.Opponent},
			card.PlayFrom{From: card.Hand, Player: card.Opponent},
		}}),
)
