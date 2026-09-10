package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
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
}
