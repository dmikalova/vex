package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hunter or Hunted?
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Choose one:
//	- Ward a creature
//	- Move a ward from a creature to another creature.
func TestHunterOrHunted(t *testing.T) {
	t.Run("wards a chosen creature", func(t *testing.T) {
		var target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				Hand:   ct.Cards(HunterOrHunted),
				InPlay: ct.Cards(ct.Bind(&target, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Play(HunterOrHunted)
		h.P1.ClickOption("ward a creature")

		if !h.Game().State.Cards[target.ID()].Warded {
			t.Error("the chosen creature should be warded")
		}
	})

	t.Run("moves a ward from one creature to another", func(t *testing.T) {
		var from, onto ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(HunterOrHunted),
				InPlay: ct.Cards(
					ct.Bind(&from, ct.Creature(ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&onto, ct.Creature(ct.Power(3))),
			)},
		})
		h.Game().State.Cards[from.ID()].Warded = true

		h.P1.Play(HunterOrHunted)
		h.P1.ClickOption("move a ward")

		if h.Game().State.Cards[from.ID()].Warded {
			t.Error("the source should lose its ward")
		}
		if !h.Game().State.Cards[onto.ID()].Warded {
			t.Error("the destination should gain the ward")
		}
	})
}
