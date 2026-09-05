package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Challe the Safeguard
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Deploy, Taunt.
func TestChalleTheSafeguard(t *testing.T) {
	t.Run("Deploy lets it enter between two creatures, and it carries Taunt", func(t *testing.T) {
		var a, b ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					ct.Bind(&a, ct.Creature(ct.Power(3))),
					ct.Bind(&b, ct.Creature(ct.Power(3))),
				),
				Hand: ct.Cards(ChalleTheSafeguard),
			},
		})

		h.P1.Play(ChalleTheSafeguard)
		h.P1.ClickOption("Between") // the sole interior option on a two-card line

		line := h.Game().Battleline(0)
		if len(line) != 3 || line[0] != a.ID() || line[2] != b.ID() {
			t.Fatalf("battleline = %v, want [a Challe b]", line)
		}
		challe := line[1]
		if h.Game().Name(challe) != "Challe the Safeguard" {
			t.Fatalf("interior card = %q, want Challe the Safeguard", h.Game().Name(challe))
		}
		if !h.Game().HasKeyword(challe, engine.Taunt) {
			t.Error("Challe the Safeguard should carry Taunt")
		}
	})
}
