package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Favor of Rex
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Trigger the play effect of a creature.
func TestFavorOfRex(t *testing.T) {
	var dux ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Saurian,
			Hand:   ct.Cards(FavorOfRex),
			InPlay: ct.Cards(ct.Bind(&dux, GrimlocusDux)),
		},
	})

	h.P1.Play(FavorOfRex) // single valid creature auto-triggers its play effect

	h.Expect(dux).AmberOn(2)
}
