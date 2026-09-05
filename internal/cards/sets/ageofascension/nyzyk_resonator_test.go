package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Nyzyk Resonator
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Martian • Soldier
//
//	For each neighbor Nyzyk Resonator has, your opponent's keys cost +2 Æmber.
func TestNyzykResonator(t *testing.T) {
	t.Run("raises the opponent's key cost by 2 per neighbor", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, InPlay: ct.Cards(
				ct.Creature(ct.Power(3)),
				NyzykResonator,
				ct.Creature(ct.Power(3)),
			)},
		})

		g := h.Game()
		g.State.Aember[1] = engine.KeyCost + 3
		g.StartTurn(1)
		h.P2.ExpectKeys(0) // one Æmber short of the +4 raised cost

		g.State.Aember[1] = engine.KeyCost + 4
		g.StartTurn(1)
		h.P2.ExpectKeys(1)
		h.P2.ExpectAmber(0) // paid the raised cost
	})

	t.Run("does not raise the key cost with no neighbors", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Mars, InPlay: ct.Cards(NyzykResonator)},
		})

		g := h.Game()
		g.State.Aember[1] = engine.KeyCost
		g.StartTurn(1)
		h.P2.ExpectKeys(1)
		h.P2.ExpectAmber(0)
	})
}
