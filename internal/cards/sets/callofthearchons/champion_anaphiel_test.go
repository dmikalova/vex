package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Champion Anaphiel
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Armor:  1
//	Traits: Knight • Spirit
//
//	Taunt.
func TestChampionAnaphiel(t *testing.T) {
	t.Run("is a 6-power creature with Taunt", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(ChampionAnaphiel),
			},
		})

		h.P1.Play(ChampionAnaphiel)

		h.Expect(ChampionAnaphiel).Power(6).At(ct.PlayArea)

		hasTaunt := false
		for _, kw := range ChampionAnaphiel.Keywords {
			if kw == card.Keyword.Taunt {
				hasTaunt = true
			}
		}
		if !hasTaunt {
			t.Errorf("Champion Anaphiel keywords = %v, want Taunt", ChampionAnaphiel.Keywords)
		}
	})
}
