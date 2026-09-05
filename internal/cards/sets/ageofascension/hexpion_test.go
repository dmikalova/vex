package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Hexpion
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Robot
//
//	Destroyed: Archive Hexpion from play. Archive the top card of your deck.
func TestHexpion(t *testing.T) {
	t.Run("archives itself and the top card of the deck when destroyed", func(t *testing.T) {
		var hexpion, foe, top ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&hexpion, Hexpion)),
				Deck:   ct.Cards(ct.Bind(&top, ct.Creature())),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6))))},
		})

		h.P1.Fight(hexpion, foe)

		h.Expect(hexpion).At(ct.Archives)
		h.Expect(top).At(ct.Archives)
	})
}
