package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// CXO Taber
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Alien • Krxix
//
//	Fight/Reap: Play or use a non-Star Alliance card.
func TestCXOTaber(t *testing.T) {
	t.Run("reaping plays one off-house card from hand", func(t *testing.T) {
		var taber ct.Card
		marsArtifact := ct.Artifact(ct.OfHouse(card.House.Mars))
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&taber, CXOTaber)),
				Hand:   ct.Cards(marsArtifact),
			},
		})

		// A non-Star Alliance card cannot be played on its own before reaping.
		h.P1.ExpectCannotPlay(marsArtifact)

		// Reaping plays the sole matching hand card at once.
		h.P1.Reap(taber)

		h.Expect(marsArtifact).At(ct.PlayArea)
	})

	t.Run("reaping uses one off-house card in play", func(t *testing.T) {
		var taber, mars ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Bind(&taber, CXOTaber),
					ct.Bind(&mars, ct.Creature(ct.OfHouse(card.House.Mars))),
				),
			},
		})

		// A Mars creature cannot be used on its own during a Star Alliance turn.
		h.P1.ExpectCannotUse(mars)

		// Reaping Taber (1 Æmber) then uses the sole Mars creature, which reaps
		// (1 more Æmber).
		h.P1.Reap(taber)

		h.Expect(mars).At(ct.PlayArea).Exhausted()
		h.P1.ExpectAmber(2)
	})
}
