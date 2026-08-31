package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// A Fair Game
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Discard the top card of your opponent's deck. Reveal your opponent's hand. For each card of the discarded card's house revealed this way, gain 1 Æmber. Discard the top card of your deck. Reveal your hand. For each card of the discarded card's house revealed this way, your opponent gains 1 Æmber.
func TestAFairGame(t *testing.T) {
	t.Run(
		"gains Æmber for each hand card matching the discarded card's house for both players",
		func(t *testing.T) {
			var opponentTop, opponentNext, yourTop, yourNext ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House: card.House.Dis,
					Hand: ct.Cards(
						AFairGame,
						ct.Tactic(ct.OfHouse(card.House.Dis)),
						ct.Tactic(ct.OfHouse(card.House.Mars)),
						ct.Creature(ct.OfHouse(card.House.Sanctum), ct.Power(2)),
					),
					Deck: ct.Cards(
						ct.Bind(&yourTop, ct.Tactic(ct.OfHouse(card.House.Dis))),
						ct.Bind(&yourNext, ct.Tactic(ct.OfHouse(card.House.Logos))),
					),
				},
				P2: ct.Side{
					Hand: ct.Cards(
						ct.Tactic(ct.OfHouse(card.House.Mars)),
						ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(2)),
						ct.Tactic(ct.OfHouse(card.House.Dis)),
					),
					Deck: ct.Cards(
						ct.Bind(&opponentTop, ct.Tactic(ct.OfHouse(card.House.Mars))),
						ct.Bind(&opponentNext, ct.Tactic(ct.OfHouse(card.House.Logos))),
					),
				},
			})

			h.P1.Play(AFairGame)

			h.Expect(opponentTop).At(ct.Discard)
			h.Expect(opponentNext).At(ct.Deck)
			h.Expect(yourTop).At(ct.Discard)
			h.Expect(yourNext).At(ct.Deck)
			h.P1.ExpectAmber(2)
			h.P2.ExpectAmber(1)
		},
	)
}
