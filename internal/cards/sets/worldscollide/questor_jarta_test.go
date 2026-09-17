package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Questor Jarta
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Dinosaur • Politician
//
//	Elusive.
//	Reap: You may exalt Questor Jarta. Gain 1 Æmber.
func TestQuestorJarta(t *testing.T) {
	t.Run("may exalt itself to gain 1 Æmber when it reaps", func(t *testing.T) {
		var jarta ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&jarta, QuestorJarta)),
			},
			P2: ct.Side{},
		})
		jarta.Ready()

		h.P1.Reap(jarta)
		h.P1.ClickCard(jarta)

		h.Expect(jarta).AmberOn(1)
		// 1 Æmber from the reap plus 1 from the gate.
		h.P1.ExpectAmber(2)
	})

	t.Run("declining takes only the reap Æmber", func(t *testing.T) {
		var jarta ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(ct.Bind(&jarta, QuestorJarta)),
			},
			P2: ct.Side{},
		})
		jarta.Ready()

		h.P1.Reap(jarta)
		h.P1.ClickDone()

		h.Expect(jarta).AmberOn(0)
		h.P1.ExpectAmber(1)
	})
}
