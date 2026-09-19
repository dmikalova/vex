package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Rhetor Gallim
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Dinosaur • Philosopher
//
//	Play: Keys cost +3 Æmber during your opponent's next turn.
//	Reap: You may exalt Rhetor Gallim. Keys cost +3 Æmber during your opponent's next turn.
func TestRhetorGallim(t *testing.T) {
	t.Run("taxes the opponent's next-turn keys when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(RhetorGallim),
			},
			P2: ct.Side{
				House: card.House.Saurian,
				Amber: 8,
			},
		})

		h.P1.Play(RhetorGallim)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Saurian)

		// 8 Æmber does not cover a key at 6 + 3.
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(8)
	})

	t.Run("may exalt itself to tax again when it reaps", func(t *testing.T) {
		var gallim ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&gallim, RhetorGallim)),
			},
			P2: ct.Side{
				House: card.House.Saurian,
				Amber: 8,
			},
		})
		gallim.Ready()

		h.P1.Reap(gallim)
		h.P1.ClickCard(gallim)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Saurian)

		h.Expect(gallim).AmberOn(1)
		h.P2.ExpectKeys(0)
		h.P2.ExpectAmber(8)
	})
}
