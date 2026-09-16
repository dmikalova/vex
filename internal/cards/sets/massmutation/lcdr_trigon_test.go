package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// LCdr. Trigon
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Mutant
//
//	Reap: Discard the top card of your deck. Resolve that card's bonus icons.
func TestLCdrTrigon(t *testing.T) {
	t.Run("reap discards the top card and resolves its bonus icons", func(t *testing.T) {
		var trigon, topCard ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&trigon, LCdrTrigon)),
				Deck: ct.Cards(
					ct.Bind(&topCard, ct.Creature(ct.Bonus(card.Bonus.Aember))),
					ct.Creature(),
				),
			},
		})

		before := h.P1.Amber()
		trigon.Ready()
		h.P1.Reap(trigon)

		// The reap gains 1 Æmber; the discarded card's Æmber bonus gains 1 more.
		if got := h.P1.Amber() - before; got != 2 {
			t.Fatalf("aember gained = %d, want 2", got)
		}
		h.Expect(topCard).At(ct.Discard)
	})
}
