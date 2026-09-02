package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nexus
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Cyborg • Thief
//
//	Elusive.
//	Reap: Use an enemy artifact as if it were yours.
func TestNexus(t *testing.T) {
	var theirs, drawn ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Shadows,
			InPlay: ct.Cards(Nexus),
			Deck:   ct.Cards(ct.Bind(&drawn, ct.Creature())),
		},
		P2: ct.Side{InPlay: ct.Cards(ct.Bind(&theirs, LibraryOfBabble))},
	})

	h.P1.Reap(Nexus)

	h.P1.ExpectAmber(1) // the reap itself
	h.Expect(drawn).At(ct.Hand)
	h.Expect(theirs).Exhausted()
}
