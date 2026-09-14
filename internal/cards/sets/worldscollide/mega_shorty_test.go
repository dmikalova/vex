package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Mega Shorty
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  6
//	Traits: Giant
//
//	Assault 4.
//	Reap: Enrage Mega Shorty.
func TestMegaShorty(t *testing.T) {
	t.Run("enrages itself when it reaps", func(t *testing.T) {
		var shorty ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&shorty, MegaShorty)),
			},
		})

		h.P1.Reap(shorty)

		if !h.Game().Enraged(shorty.ID()) {
			t.Errorf("%s should be enraged", shorty.Name())
		}
	})
}
