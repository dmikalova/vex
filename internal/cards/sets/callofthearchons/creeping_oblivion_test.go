package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Creeping Oblivion
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Purge up to 2 cards from a discard pile.
func TestCreepingOblivion(t *testing.T) {
	t.Run("purges up to two cards from a discard pile", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(CreepingOblivion)},
			P2: ct.Side{Discard: ct.Cards(
				ct.Bind(&a, ct.Creature(ct.OfHouse(card.House.Brobnar))),
				ct.Bind(&b, ct.Creature(ct.OfHouse(card.House.Brobnar))),
			)},
		})

		h.P1.Play(CreepingOblivion)
		h.P1.ClickOption(a.Name())
		h.P1.ClickOption(b.Name())

		h.Expect(a).At(ct.Purge)
		h.Expect(b).At(ct.Purge)
	})
}
