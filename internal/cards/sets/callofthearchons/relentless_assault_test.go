package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Relentless Assault
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Ready and fight with up to 3 different friendly creatures, one at a time.
func TestRelentlessAssault(t *testing.T) {
	brobnar := ct.OfHouse(card.House.Brobnar)

	setup := func(t *testing.T) (h *ct.Harness, troll, giant, witch ct.Card) {
		h = ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				Hand:  ct.Cards(RelentlessAssault),
				InPlay: ct.Cards(
					ct.Bind(&troll, ct.Creature(brobnar, ct.Power(2))),
					ct.Bind(&giant, ct.Creature(brobnar, ct.Power(3))),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&witch, ct.Creature(brobnar, ct.Power(20))),
					ct.Creature(
						brobnar,
						ct.Power(20),
					), // a second enemy, so the fight is a real choice
				),
			},
		})
		return h, troll, giant, witch
	}

	t.Run("fights with each chosen creature in turn", func(t *testing.T) {
		h, troll, giant, witch := setup(t)

		h.P1.Play(RelentlessAssault)
		h.P1.ExpectPrompt("Choose a friendly creature").Source("Relentless Assault")
		h.P1.ClickCard(troll)
		h.P1.ClickCard(witch)
		h.P1.ClickCard(giant)
		h.P1.ClickCard(witch)

		// Both traded into the 20-power witch and died; she took 2 + 3 damage.
		h.Expect(troll).At(ct.Discard)
		h.Expect(giant).At(ct.Discard)
		h.Expect(witch).At(ct.PlayArea).Damage(5)
	})

	t.Run("stops early when the controller declines a pass", func(t *testing.T) {
		h, troll, giant, witch := setup(t)

		h.P1.Play(RelentlessAssault)
		h.P1.ClickCard(troll)
		h.P1.ClickCard(witch)
		h.P1.ClickDone()

		h.Expect(troll).At(ct.Discard)
		h.Expect(giant).At(ct.PlayArea).Ready()
		h.Expect(witch).At(ct.PlayArea).Damage(2)
	})
}
