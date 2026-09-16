package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Subject Kirby
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Play/Fight/Reap: You may play a non-Star Alliance creature this turn.
func TestSubjectKirby(t *testing.T) {
	t.Run("lets a non-Star Alliance creature be played this turn", func(t *testing.T) {
		var kirby, offHouse ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&kirby, SubjectKirby)),
				Hand:   ct.Cards(ct.Bind(&offHouse, ct.Creature(ct.OfHouse(card.House.Untamed)))),
			},
		})

		h.P1.Reap(kirby)
		h.P1.Play(offHouse)

		h.Expect(offHouse).At(ct.PlayArea)
	})
}
