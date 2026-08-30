package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Raiding Knight
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Play: Raiding Knight captures 1 Æmber from your opponent.
func TestRaidingKnight(t *testing.T) {
	t.Run("captures 1 Æmber when played", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, Hand: ct.Cards(RaidingKnight)},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Play(RaidingKnight)

		h.Expect(RaidingKnight).AmberOn(1)
		h.P2.ExpectAmber(2)
	})
}
