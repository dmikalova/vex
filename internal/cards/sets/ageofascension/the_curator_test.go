package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// The Curator
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Scientist
//
//	Friendly artifacts enter play ready.
func TestTheCurator(t *testing.T) {
	var art ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Logos,
			InPlay: ct.Cards(TheCurator),
			Hand:   ct.Cards(ct.Bind(&art, ct.Artifact(ct.OfHouse(card.House.Logos)))),
		},
	})

	h.P1.Play(art)

	h.Expect(art).Ready()
}
