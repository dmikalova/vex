package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Musthic Murmook
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast
//
//	Each player's keys cost +1 Æmber.
//	Play: Deal 4 damage to a Creature.
func TestMusthicMurmook(t *testing.T) {
	t.Run("deals 4 damage to a creature when played", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(MusthicMurmook)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(10))))},
		})

		h.P1.Play(MusthicMurmook)
		h.P1.ClickCard(foe)

		h.Expect(foe).Damage(4)
	})

	t.Run("raises each player's key cost by 1 while in play", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, InPlay: ct.Cards(MusthicMurmook)},
		})

		g := h.Game()
		g.State.Aember[0] = engine.KeyCost
		g.StartTurn(0)
		h.P1.ExpectKeys(0) // one Æmber short of the raised cost

		g.State.Aember[0] = engine.KeyCost + 1
		g.StartTurn(0)
		h.P1.ExpectKeys(1)
	})
}
