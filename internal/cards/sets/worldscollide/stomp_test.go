package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Stomp
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Deal 5 damage to a creature. If this damage destroys that creature, exalt a friendly creature.
func TestStomp(t *testing.T) {
	t.Run("exalts a friendly creature when the damage destroys the target", func(t *testing.T) {
		var foe, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(Stomp),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(4)))),
			},
		})

		h.P1.Play(Stomp)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.Discard)
		h.Expect(ally).AmberOn(1)
	})

	t.Run("does not exalt when the target survives", func(t *testing.T) {
		var foe, ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(Stomp),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.OfHouse(card.House.Saurian)))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(8)))),
			},
		})

		h.P1.Play(Stomp)
		h.P1.ClickCard(foe)

		h.Expect(foe).At(ct.PlayArea).Damage(5)
		h.Expect(ally).AmberOn(0)
	})
}
