package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// The Archivist
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Cyborg
//
//	Instead of picking up all of your archives, you may pick up any number of cards in your archives.
func TestTheArchivist(t *testing.T) {
	var taken, kept ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Logos,
			InPlay: ct.Cards(TheArchivist),
			Archives: ct.Cards(
				ct.Bind(&taken, ct.Creature(ct.Power(2))),
				ct.Bind(&kept, ct.Creature(ct.Power(2))),
			),
		},
	})

	// On a fresh turn, choosing a house offers the archives. With The Archivist in
	// play the pickup is selective: take one card into hand and leave the other.
	h.P1.EndTurn()
	h.P2.EndTurn()
	h.P1.ChooseHouse(card.House.Logos)
	h.P1.ClickCard(taken)
	h.P1.ClickDone()

	h.Expect(taken).At(ct.Hand)
	h.Expect(kept).At(ct.Archives)
}
