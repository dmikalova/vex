package worldscollide

import (
	"slices"
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Gebuk
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Beast
//
//	Destroyed: Discard the top card of your deck. If it is a creature, swap it with Gebuk.
func TestGebuk(t *testing.T) {
	t.Run("a creature off the top of the deck enters play in Gebuk's slot", func(t *testing.T) {
		var gebuk, left, right, reborn ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.Power(4))),
					ct.Bind(&gebuk, Gebuk),
					ct.Bind(&right, ct.Creature(ct.Power(4))),
				),
				Deck: ct.Cards(ct.Bind(&reborn, ct.Creature(ct.Power(3)))),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{gebuk.ID()})

		h.Expect(gebuk).At(ct.Discard)
		h.Expect(reborn).Exhausted().At(ct.PlayArea)
		if want := []engine.LocalID{
			left.ID(),
			reborn.ID(),
			right.ID(),
		}; !slices.Equal(
			h.Game().Battleline(0),
			want,
		) {
			t.Errorf(
				"battleline = %v, want %v (reborn in Gebuk's slot)",
				h.Game().Battleline(0),
				want,
			)
		}
	})

	t.Run("a non-creature off the top is discarded and nothing enters play", func(t *testing.T) {
		var gebuk, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Untamed,
				InPlay: ct.Cards(ct.Bind(&gebuk, Gebuk)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Tactic())),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{gebuk.ID()})

		h.Expect(gebuk).At(ct.Discard)
		h.Expect(top).At(ct.Discard)
		if n := len(h.Game().Battleline(0)); n != 0 {
			t.Errorf("battleline has %d creatures, want 0", n)
		}
	})
}
