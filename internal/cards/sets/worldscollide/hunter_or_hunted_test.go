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
//	Play: Remove a ward from a creature, and ward a creature.
func TestHunterOrHunted(t *testing.T) {
	t.Run("removes a ward from one creature and wards another", func(t *testing.T) {
		var warded, target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(HunterOrHunted),
				InPlay: ct.Cards(
					ct.Bind(&warded, ct.Creature(ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&target, ct.Creature(ct.Power(3))),
			)},
		})
		h.Game().State.Cards[warded.ID()].Warded = true

		h.P1.Play(HunterOrHunted)
		h.P1.ClickCard(warded) // remove the ward from this creature
		h.P1.ClickCard(target) // ward this creature

		if h.Game().State.Cards[warded.ID()].Warded {
			t.Error("the chosen creature should lose its ward")
		}
		if !h.Game().State.Cards[target.ID()].Warded {
			t.Error("the second chosen creature should be warded")
		}
	})

	t.Run("wards a creature when there is no ward to remove", func(t *testing.T) {
		var plain, target ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Shadows,
				Hand:  ct.Cards(HunterOrHunted),
				InPlay: ct.Cards(
					ct.Bind(&plain, ct.Creature(ct.Power(3))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&target, ct.Creature(ct.Power(3))),
			)},
		})

		h.P1.Play(HunterOrHunted)
		h.P1.ClickCard(plain)  // no ward to remove — nothing happens
		h.P1.ClickCard(target) // ward this creature

		if h.Game().State.Cards[plain.ID()].Warded {
			t.Error("removing from an unwarded creature should leave it unwarded")
		}
		if !h.Game().State.Cards[target.ID()].Warded {
			t.Error("the second chosen creature should be warded")
		}
	})
}
