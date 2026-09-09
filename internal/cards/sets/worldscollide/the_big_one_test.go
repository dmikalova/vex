package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Big One
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Weapon
//
//	After a creature is played, put a fuse counter on The Big One. If there are 10 or more fuse counters on The Big One, destroy each creature and each artifact.
func TestTheBigOne(t *testing.T) {
	t.Run("puts a fuse counter on itself after a creature is played", func(t *testing.T) {
		var bomb, newbie ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				Hand:   ct.Cards(ct.Bind(&newbie, ct.Creature(ct.Power(4)))),
				InPlay: ct.Cards(ct.Bind(&bomb, TheBigOne)),
			},
		})

		h.P1.Play(newbie)

		if got := h.Game().CountersOn(bomb.ID(), card.Counter.Fuse); got != 1 {
			t.Errorf("fuse counters = %d, want 1", got)
		}
	})

	t.Run("destroys each creature and artifact once ten fuse counters accrue", func(t *testing.T) {
		var bomb, newbie, other ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(ct.Bind(&newbie, ct.Creature(ct.Power(4)))),
				InPlay: ct.Cards(
					ct.Bind(&bomb, TheBigOne),
					ct.Bind(&other, ct.Creature(ct.Power(4))),
				),
			},
		})
		// Nine fuse counters already sit on The Big One; the tenth is the trigger.
		h.Game().PlaceCounter(bomb.ID(), card.Counter.Fuse, 9)

		h.P1.Play(newbie)

		h.Expect(bomb).At(ct.Discard)
		h.Expect(other).At(ct.Discard)
		h.Expect(newbie).At(ct.Discard)
	})
}
