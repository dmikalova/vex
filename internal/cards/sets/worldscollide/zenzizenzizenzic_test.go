package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Zenzizenzizenzic
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Cyborg • Leader
//
//	While Zenzizenzizenzic is in the center of the battleline, your hand size is 2 more.
func TestZenzizenzizenzic(t *testing.T) {
	t.Run("refills two extra while in the center", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(Zenzizenzizenzic),
				Deck:   ct.DeckOf(card.House.Logos, 12),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()

		if got := int(h.Game().State.Hand[0].Count); got != 8 {
			t.Errorf("hand after draw = %d, want 8", got)
		}
	})

	t.Run("does nothing while off the center", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					Zenzizenzizenzic,
					ct.Creature(ct.Power(3)),
				),
				Deck: ct.DeckOf(card.House.Logos, 12),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()

		if got := int(h.Game().State.Hand[0].Count); got != 6 {
			t.Errorf("hand after draw = %d, want 6", got)
		}
	})
}
