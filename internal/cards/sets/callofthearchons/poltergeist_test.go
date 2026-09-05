package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Poltergeist
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Use an artifact. Destroy it.
func TestPoltergeist(t *testing.T) {
	t.Run("uses an enemy artifact and then destroys it", func(t *testing.T) {
		var theirs, drawn ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(Poltergeist),
				Deck:  ct.Cards(ct.Bind(&drawn, ct.Creature())),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&theirs, LibraryOfBabble))},
		})

		h.P1.Play(Poltergeist)

		h.Expect(drawn).At(ct.Hand)
		h.Expect(theirs).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})

	t.Run("can use and destroy your own artifact", func(t *testing.T) {
		var mine, drawn ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				Hand:   ct.Cards(Poltergeist),
				Deck:   ct.Cards(ct.Bind(&drawn, ct.Creature())),
				InPlay: ct.Cards(ct.Bind(&mine, LibraryOfBabble)),
			},
		})

		h.P1.Play(Poltergeist)

		h.Expect(drawn).At(ct.Hand)
		h.Expect(mine).At(ct.Discard)
	})
}
