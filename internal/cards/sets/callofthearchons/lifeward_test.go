package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Lifeward
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Power
//
//	Versatile.
//	Action: Destroy Lifeward. Your opponent cannot play creatures during their next turn.
func TestLifeward(t *testing.T) {
	t.Run(
		"destroys itself and bars the opponent from playing creatures next turn",
		func(t *testing.T) {
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{House: card.House.Dis, InPlay: ct.Cards(Lifeward)},
			})

			h.P1.UseAction(Lifeward)

			h.Expect(Lifeward).At(ct.Discard)
			if got := h.Game().State.CannotPlayTypeNext[1]; got != engine.Creature {
				t.Errorf("opponent's armed play bar = %q, want Creature", got)
			}
		},
	)
}
