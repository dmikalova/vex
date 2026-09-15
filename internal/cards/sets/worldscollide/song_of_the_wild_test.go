package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Song of the Wild
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Each friendly creature gains, "Reap: Gain 1 Æmber."
func TestSongOfTheWild(t *testing.T) {
	t.Run("friendly creatures gain 1 Æmber when they reap this turn", func(t *testing.T) {
		var beast ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				Hand:   ct.Cards(SongOfTheWild),
				InPlay: ct.Cards(ct.Bind(&beast, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Play(SongOfTheWild)
		h.P1.Reap(beast)

		// 1 Æmber from the reap itself, 1 more from the granted ability.
		h.P1.ExpectAmber(2)
	})
}
