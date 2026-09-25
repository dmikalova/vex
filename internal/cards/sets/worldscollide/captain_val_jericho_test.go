package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Captain Val Jericho
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Human • Leader
//
//	During your turn, if Captain Val Jericho is in the center of your battleline, you may play one card that is not of the active house.
func TestCaptainValJericho(t *testing.T) {
	t.Run("frees one non-active play while centered", func(t *testing.T) {
		var jericho ct.Card
		marsA := ct.Creature(ct.OfHouse(card.House.Mars))
		marsB := ct.Creature(ct.OfHouse(card.House.Mars))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&jericho, CaptainValJericho)),
				Hand:   ct.Cards(marsA, marsB),
			},
		})

		// The sole creature is centered, so one off-house play is freed; the
		// second is not.
		h.P1.Play(marsA)
		h.P1.ExpectCannotPlay(marsB)
	})

	// A grant can modify the first-turn rule: Val Jericho is P1's one first-turn
	// play, and the off-house play it frees is not barred by the one-card limit.
	t.Run("frees a further play on the first player's first turn", func(t *testing.T) {
		var mars ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				Hand: ct.Cards(
					CaptainValJericho,
					ct.Bind(&mars, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})
		h.Game().State.FirstTurnPlayLimit[0] = true

		h.P1.Play(CaptainValJericho)
		h.P1.Play(mars)
		h.Expect(mars).At(ct.PlayArea)
	})
}
