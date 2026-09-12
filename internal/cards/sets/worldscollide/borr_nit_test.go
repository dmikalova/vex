package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Borr Nit
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Demon
//
//	Reap: Reveal the top 5 cards of a player's deck. Purge a card revealed this way. Shuffle that deck.
func TestBorrNit(t *testing.T) {
	var borr, victim, keep ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Dis,
			InPlay: ct.Cards(ct.Bind(&borr, BorrNit)),
		},
		P2: ct.Side{
			Deck: ct.Cards(
				ct.Bind(&victim, ct.Creature()),
				ct.Bind(&keep, ct.Creature()),
			),
		},
	})

	h.P1.Reap(borr)
	h.P1.ClickOption("opponent") // reveal from the opponent's deck
	h.P1.ClickCard(victim)       // purge the revealed victim

	h.Expect(victim).At(ct.Purge)
	if deck := h.Game().Deck(1); len(deck) != 1 || deck[0] != keep.ID() {
		t.Errorf("opponent deck = %v, want [%d]", deck, keep.ID())
	}
}
