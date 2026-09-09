package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// The Red Baron
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  4
//	Armor:  1
//	Traits: Cyborg • Pirate
//
//	While your opponent's red key is forged, The Red Baron gains elusive.
//	While your red key is forged, The Red Baron gains, "Reap: Steal 1 Æmber."
func TestTheRedBaron(t *testing.T) {
	t.Run("reaps to steal while your red key is forged", func(t *testing.T) {
		var baron ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:      card.House.Brobnar,
				ForgedKeys: []engine.KeyColor{card.KeyColor.Red},
				InPlay:     ct.Cards(ct.Bind(&baron, TheRedBaron)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(baron)

		// Base reap gains 1, the granted reap steals 1 more from the opponent.
		h.P1.ExpectAmber(2)
		h.P2.ExpectAmber(2)
	})

	t.Run("no steal while your red key is not forged", func(t *testing.T) {
		var baron ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&baron, TheRedBaron)),
			},
			P2: ct.Side{Amber: 3},
		})

		h.P1.Reap(baron)

		h.P1.ExpectAmber(1)
		h.P2.ExpectAmber(3)
	})

	t.Run("gains elusive while your opponent's red key is forged", func(t *testing.T) {
		var baron ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&baron, TheRedBaron)),
			},
			P2: ct.Side{ForgedKeys: []engine.KeyColor{card.KeyColor.Red}},
		})

		if !h.Game().HasKeyword(baron.ID(), card.Keyword.Elusive) {
			t.Fatal("The Red Baron should have elusive while the opponent's red key is forged")
		}
	})

	t.Run("no elusive while your opponent's red key is not forged", func(t *testing.T) {
		var baron ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&baron, TheRedBaron)),
			},
		})

		if h.Game().HasKeyword(baron.ID(), card.Keyword.Elusive) {
			t.Fatal("The Red Baron should not have elusive while no red key is forged")
		}
	})
}
