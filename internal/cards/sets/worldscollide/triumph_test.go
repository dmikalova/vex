package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Triumph
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If there are no enemy creatures in play, exalt each friendly creature. If there are 6 or more friendly creatures in play, forge a key at no cost.
func TestTriumph(t *testing.T) {
	t.Run("with no enemy creatures and 6 friendly, exalts each and forges", func(t *testing.T) {
		var watched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Saurian,
				Hand:  ct.Cards(Triumph),
				InPlay: ct.Cards(
					ct.Bind(&watched, ct.Creature()),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
					ct.Creature(),
				),
			},
		})

		h.P1.Play(Triumph)

		h.Expect(watched).AmberOn(1)
		h.P1.ExpectKeys(1)
	})

	t.Run("with fewer than 6 friendly, exalts each but does not forge", func(t *testing.T) {
		var watched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(Triumph),
				InPlay: ct.Cards(ct.Bind(&watched, ct.Creature()), ct.Creature()),
			},
		})

		h.P1.Play(Triumph)

		h.Expect(watched).AmberOn(1)
		h.P1.ExpectKeys(0)
	})

	t.Run("with enemy creatures, does nothing", func(t *testing.T) {
		var watched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(Triumph),
				InPlay: ct.Cards(ct.Bind(&watched, ct.Creature())),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Creature())},
		})

		h.P1.Play(Triumph)

		h.Expect(watched).AmberOn(0)
		h.P1.ExpectKeys(0)
	})
}
