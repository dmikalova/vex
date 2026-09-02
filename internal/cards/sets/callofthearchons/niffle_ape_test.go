package callofthearchons_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	cota "github.com/dmikalova/vactrol/internal/cards/sets/callofthearchons"
)

// Niffle Ape
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast • Niffle
//
//	While Niffle Ape is attacking, ignore taunt and elusive.
func TestNiffleApe(t *testing.T) {
	var ape, shielded, elusive ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:  card.House.Untamed,
			InPlay: ct.Cards(ct.Bind(&ape, cota.NiffleApe)),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Creature(ct.Power(4), ct.Keywords(card.Keyword.Taunt)),
				ct.Bind(&shielded, ct.Creature(ct.Power(1))),
				ct.Bind(&elusive, ct.Creature(ct.Power(1), ct.Keywords(card.Keyword.Elusive))),
			),
		},
	})

	// Taunt does not shield its neighbor from the ape, and elusive does not stop the
	// damage, so a 1-power creature dies to the first swing.
	h.P1.Fight(ape, shielded)
	h.Expect(shielded).At(ct.Discard)

	h.P1.EndTurn()
	h.P2.EndTurn()
	h.P1.ChooseHouse(card.House.Untamed)
	h.P1.Fight(ape, elusive)
	h.Expect(elusive).At(ct.Discard)
}
