package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Noname
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Demon
//
//	Noname gets +1 power for each purged card.
//	Play/Fight/Reap: Purge a card in a discard pile.
func TestNoname(t *testing.T) {
	t.Run("reaping purges a discard-pile card and grows Noname", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Dis,
				InPlay:  ct.Cards(Noname),
				Discard: ct.Cards(ct.Tactic(ct.OfHouse(card.House.Dis))),
			},
			P2: ct.Side{},
		})

		h.Expect(Noname).Power(1)

		h.P1.Reap(Noname)

		h.Expect(Noname).Power(2)
	})
}
