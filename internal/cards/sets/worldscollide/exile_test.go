package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Exile
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Your opponent gains control of a friendly creature.
func TestExile(t *testing.T) {
	t.Run("gives control of a friendly creature to the opponent", func(t *testing.T) {
		var pet ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				Hand:   ct.Cards(Exile),
				InPlay: ct.Cards(ct.Bind(&pet, ct.Creature(ct.OfHouse(card.House.Saurian)))),
			},
		})

		h.P1.Play(Exile)

		h.Expect(pet).At(ct.PlayArea)
		if got := h.Game().Controller(pet.ID()); got != 1 {
			t.Errorf("controller = %d, want 1 (opponent)", got)
		}
	})
}
