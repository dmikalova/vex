package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Remote Access
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Use an enemy artifact as if it were yours.
func TestRemoteAccess(t *testing.T) {
	t.Run("uses an enemy artifact, resolving it for you", func(t *testing.T) {
		var theirs, drawn ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				Hand:  ct.Cards(RemoteAccess),
				Deck:  ct.Cards(ct.Bind(&drawn, ct.Creature())),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&theirs, LibraryOfBabble))},
		})

		h.P1.Play(RemoteAccess)

		h.Expect(drawn).At(ct.Hand) // the draw went to the user, not the owner
		h.Expect(theirs).Exhausted()
		h.P1.ExpectAmber(1)
	})

	t.Run("does nothing when the opponent controls no artifact", func(t *testing.T) {
		var mine ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				Hand:   ct.Cards(RemoteAccess),
				InPlay: ct.Cards(ct.Bind(&mine, LibraryOfBabble)),
			},
		})

		h.P1.Play(RemoteAccess)

		h.Expect(mine).Ready() // only enemy artifacts are offered
		h.P1.ExpectAmber(1)
	})
}
