package callofthearchons

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Transposition Sandals
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Action: Swap this creature with another friendly creature in your battleline. Use the other creature."
func TestTranspositionSandals(t *testing.T) {
	t.Run("swaps the host with another friendly creature and uses that creature", func(t *testing.T) {
		var left, host, other, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Bind(&left, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
					ct.Upgraded(ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))), TranspositionSandals),
					ct.Bind(&other, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(2))),
					ct.Bind(&right, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
				),
			},
		})

		h.P1.UseAction(host)
		h.P1.ExpectPrompt("Choose another friendly creature")
		h.P1.ClickCard(other)

		want := []engine.LocalID{left.ID(), other.ID(), host.ID(), right.ID()}
		if got := h.Game().Battleline(0); !slices.Equal(got, want) {
			t.Fatalf("battleline = %v, want %v", got, want)
		}
		h.Expect(other).Exhausted()
		h.P1.ExpectAmber(1)
	})
}
