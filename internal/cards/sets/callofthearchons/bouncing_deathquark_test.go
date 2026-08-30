package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Bouncing Deathquark
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy an enemy creature and a friendly creature -> if there is a friendly creature in play, you may repeat this effect.
func TestBouncingDeathquark(t *testing.T) {
	t.Run("destroys one enemy and one friendly creature", func(t *testing.T) {
		var friend, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&friend, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)))),
				Hand:   ct.Cards(BouncingDeathquark),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3)))),
			},
		})

		// One creature on each side: both destructions auto-resolve, and with no
		// friendly creature left the repeat is never offered.
		h.P1.Play(BouncingDeathquark)

		h.Expect(friend).At(ct.Discard)
		h.Expect(foe).At(ct.Discard)
	})

	t.Run("may repeat while a friendly creature remains", func(t *testing.T) {
		var friendA, friendB, foeA, foeB ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&friendA, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
					ct.Bind(&friendB, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
				Hand: ct.Cards(BouncingDeathquark),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&foeA, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
					ct.Bind(&foeB, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(3))),
				),
			},
		})

		h.P1.Play(BouncingDeathquark)
		h.P1.ClickCard(foeA)    // first enemy to destroy
		h.P1.ClickCard(friendA) // first friendly to destroy
		h.P1.ExpectPrompt("Repeat this effect?").Source("Bouncing Deathquark")
		h.P1.ClickOption("Yes") // a friendly remains: repeat, destroying the rest

		h.Expect(foeA).At(ct.Discard)
		h.Expect(friendA).At(ct.Discard)
		h.Expect(foeB).At(ct.Discard)
		h.Expect(friendB).At(ct.Discard)
	})
}
