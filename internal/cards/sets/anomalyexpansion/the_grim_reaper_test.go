package anomalyexpansion

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Grim Reaper
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Traits: Robot • Specter
//
//	If you are haunted, The Grim Reaper enters play ready.
//	Reap: Purge an enemy creature, and purge a friendly creature.
func TestTheGrimReaper(t *testing.T) {
	t.Run("enters play ready while haunted", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Brobnar,
				Hand:    ct.Cards(ct.Bind(&reaper, TheGrimReaper)),
				Discard: ct.DeckOf(card.House.Brobnar, 10),
			},
		})

		h.P1.Play(TheGrimReaper)

		h.Expect(reaper).Ready()
	})

	t.Run("enters play exhausted when not haunted", func(t *testing.T) {
		var reaper ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Brobnar,
				Hand:    ct.Cards(ct.Bind(&reaper, TheGrimReaper)),
				Discard: ct.DeckOf(card.House.Brobnar, 9),
			},
		})

		h.P1.Play(TheGrimReaper)

		h.Expect(reaper).Exhausted()
	})

	t.Run("reap purges one enemy and one friendly creature", func(t *testing.T) {
		var reaper, enemy, friend ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&reaper, TheGrimReaper),
					ct.Bind(&friend, ct.Creature()),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&enemy, ct.Creature())),
			},
		})

		h.P1.Reap(reaper)
		// The lone enemy auto-resolves the enemy purge; the friendly purge
		// chooses between the reaper and its friend.
		h.P1.ClickCard(friend) // purge a friendly creature

		h.Expect(enemy).At(ct.Purge)
		h.Expect(friend).At(ct.Purge)
	})
}
