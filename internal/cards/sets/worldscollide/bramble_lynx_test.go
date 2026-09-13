package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bramble Lynx
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast
//
//	Skirmish.
//	If you have used a Creature to reap this turn, Bramble Lynx enters play ready.
func TestBrambleLynx(t *testing.T) {
	t.Run("enters play ready once you have reaped this turn", func(t *testing.T) {
		var lynx, reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ct.Bind(&lynx, BrambleLynx)),
				InPlay: ct.Cards(
					ct.Bind(&reaper, ct.Creature(ct.OfHouse(card.House.Untamed), ct.Power(3))),
				),
			},
		})

		h.P1.Reap(reaper)
		h.P1.Play(BrambleLynx)

		if lynx.Exhausted() {
			t.Error("Bramble Lynx should enter play ready after a reap")
		}
	})

	t.Run("enters play exhausted when you have not reaped", func(t *testing.T) {
		var lynx ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Untamed, Hand: ct.Cards(ct.Bind(&lynx, BrambleLynx))},
		})

		h.P1.Play(BrambleLynx)

		if !lynx.Exhausted() {
			t.Error("Bramble Lynx should enter play exhausted with no reap")
		}
	})
}
