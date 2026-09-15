package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Hapless Cadet
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Alien
//
//	Taunt.
//	Destroyed: Your opponent loses 3 Æmber.
func TestHaplessCadet(t *testing.T) {
	t.Run("opponent loses 3 Æmber when it is destroyed", func(t *testing.T) {
		var cadet ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.StarAlliance,
				InPlay: ct.Cards(ct.Bind(&cadet, HaplessCadet)),
			},
			P2: ct.Side{Amber: 5},
		})

		h.Game().DestroyEach(0, []engine.LocalID{cadet.ID()})

		h.P2.ExpectAmber(2)
	})
}
