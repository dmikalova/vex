package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Uxlyx the Zookeeper
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Martian • Scientist
//
//	Elusive.
//	Reap: Put an enemy creature into your archives.
func TestUxlyxTheZookeeper(t *testing.T) {
	t.Run("abducts an enemy creature, which goes home when the archives empty", func(t *testing.T) {
		var uxlyx, prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Mars,
				InPlay: ct.Cards(ct.Bind(&uxlyx, UxlyxTheZookeeper)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&prey, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Reap(uxlyx)
		h.Expect(prey).At(ct.Archives)

		h.P1.EndTurn()
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Mars)
		h.P1.ClickOption("Yes")

		h.Expect(prey).At(ct.Hand)
		if got := h.Game().Owner(prey.ID()); got != 1 {
			t.Errorf("owner = %d, want 1", got)
		}
		for _, id := range h.Game().Hand(0) {
			if id == prey.ID() {
				t.Error("the abducted creature should go to its owner's hand, not the abductor's")
			}
		}
	})
}
