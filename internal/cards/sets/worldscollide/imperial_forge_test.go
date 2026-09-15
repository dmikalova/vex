package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Imperial Forge
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Forge a key at +8 Æmber current cost, reduced by 1 Æmber for each Æmber on friendly creatures -> purge Imperial Forge.
func TestImperialForge(t *testing.T) {
	var creature ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Saurian,
			Hand:  ct.Cards(ImperialForge),
			InPlay: ct.Cards(
				ct.Bind(&creature, ct.Creature(ct.OfHouse(card.House.Saurian))),
			),
			Amber: 10,
		},
	})
	// 6 Æmber sit on the friendly creature, cutting the +8 surcharge to +2, so a
	// key that normally costs 6 costs 8 here. Imperial Forge's own 1 Æmber bonus
	// tops the pool to 11 as it is played, leaving 3 after the 8 is paid.
	h.Game().State.Cards[creature.ID()].Amber = 6

	h.P1.Play(ImperialForge)

	if got := h.Game().Keys(0); got != 1 {
		t.Errorf("keys = %d, want 1", got)
	}
	h.P1.ExpectAmber(3)
}
