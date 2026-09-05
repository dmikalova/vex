package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Skullion
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  7
//	Armor:  2
//	Traits: Demon
//
//	Play: Destroy a friendly creature.
func TestSkullion(t *testing.T) {
	t.Run("sacrifices a chosen friendly creature", func(t *testing.T) {
		var ally ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				Hand:   ct.Cards(Skullion),
				InPlay: ct.Cards(ct.Bind(&ally, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Play(Skullion)
		h.P1.ClickCard(ally)

		h.Expect(ally).At(ct.Discard)
	})
}
