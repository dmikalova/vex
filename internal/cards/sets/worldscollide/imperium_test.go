package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Imperium
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Ward 2 friendly creatures.
func TestImperium(t *testing.T) {
	var a, b, c ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			Hand:  ct.Cards(ct.Bind(&a, Imperium)),
			InPlay: ct.Cards(
				ct.Bind(&b, ct.Creature()),
				ct.Bind(&c, ct.Creature()),
				ct.Creature(),
			),
		},
	})

	h.P1.Play(a)
	h.P1.ClickCard(b)
	h.P1.ClickCard(c)

	if !h.Game().Warded(b.ID()) {
		t.Errorf("creature b not warded")
	}
	if !h.Game().Warded(c.ID()) {
		t.Errorf("creature c not warded")
	}
}
