package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Phalanx Strike
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For each friendly creature in play, deal 1 damage to a creature. You may exalt a friendly creature to repeat the preceding effect.
func TestPhalanxStrike(t *testing.T) {
	t.Run("deals 1 damage per friendly creature, then declines the exalt", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(PhalanxStrike),
				InPlay: ct.Cards(
					ct.Creature(ct.OfHouse(card.House.Saurian)),
					ct.Creature(ct.OfHouse(card.House.Saurian)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(9))))},
		})

		h.P1.Play(PhalanxStrike)
		h.P1.ClickCard(enemy)
		h.P1.ClickDone() // decline the exalt

		h.Expect(enemy).Damage(2)
	})

	t.Run("exalting a friendly creature repeats the effect", func(t *testing.T) {
		var enemy, pay ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(PhalanxStrike),
				InPlay: ct.Cards(
					ct.Bind(&pay, ct.Creature(ct.OfHouse(card.House.Saurian))),
					ct.Creature(ct.OfHouse(card.House.Saurian)),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature(ct.Power(9))))},
		})

		h.P1.Play(PhalanxStrike)
		h.P1.ClickCard(enemy) // first hit: 2 damage
		h.P1.ClickCard(pay)   // exalt pay to repeat
		h.P1.ClickCard(enemy) // second hit: 2 more damage
		h.P1.ClickDone()      // decline further exalts

		h.Expect(enemy).Damage(4)
		if got := h.Game().State.Cards[pay.ID()].Amber; got != 1 {
			t.Errorf("exalted amber = %d, want 1", got)
		}
	})
}
