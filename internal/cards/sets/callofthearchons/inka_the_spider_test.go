package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Inka the Spider
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast
//
//	Poison.
//	Play/Reap: Stun a creature.
func TestInkaTheSpider(t *testing.T) {
	t.Run("stuns a chosen creature when played", func(t *testing.T) {
		var foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(InkaTheSpider)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars))))},
		})

		h.P1.Play(InkaTheSpider)
		h.P1.ClickCard(foe)

		h.Expect(foe).Stunned(true)
	})
}
