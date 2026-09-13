package ageofascension_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/cards/sets/ageofascension"
)

// Anahita the Trader
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Merchant
//
//	Reap: Your opponent gains control of a friendly Artifact -> steal 2 Æmber.
func TestAnahitaTheTrader(t *testing.T) {
	t.Run("gives an artifact away and takes 2 Æmber for it", func(t *testing.T) {
		var anahita, relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Amber: 0,
				InPlay: ct.Cards(
					ct.Bind(&anahita, ageofascension.AnahitaTheTrader),
					ct.Bind(&relic, ct.Artifact()),
				),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(anahita)

		// 1 Æmber from reaping, plus 2 stolen from the opponent for the artifact.
		h.P1.ExpectAmber(3)
		h.P2.ExpectAmber(1)

		g := h.Game()
		found := false
		for _, id := range g.Artifacts(1) {
			if g.Name(id) == g.Name(relic.ID()) {
				found = true
			}
		}
		if !found {
			t.Error("the artifact should now be controlled by the opponent")
		}
	})

	t.Run("does nothing when there is no friendly artifact to give", func(t *testing.T) {
		var anahita ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Sanctum,
				Amber:  0,
				InPlay: ct.Cards(ct.Bind(&anahita, ageofascension.AnahitaTheTrader)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(anahita)

		// Only the base reap Æmber; the gate never fires without an artifact.
		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(3)
	})
}
