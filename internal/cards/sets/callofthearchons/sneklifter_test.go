package callofthearchons

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Sneklifter
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: Take control of an enemy artifact. If it does not belong to a house on your identity, it belongs to house Shadows until it leaves play.
func TestSneklifter(t *testing.T) {
	t.Run(
		"takes control of an enemy artifact and reassigns an off-identity one to Shadows",
		func(t *testing.T) {
			var relic ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(Sneklifter)},
				P2: ct.Side{
					InPlay: ct.Cards(ct.Bind(&relic, ct.Artifact(ct.OfHouse(card.House.Untamed)))),
				},
			})
			h.Game().
				SetPlayerHouses(0, []engine.House{card.House.Shadows, card.House.Logos, card.House.Sanctum})

			h.P1.Play(Sneklifter)

			if !slices.Contains(h.Game().Artifacts(0), relic.ID()) {
				t.Error("the artifact should be controlled by P1")
			}
			if got := h.Game().House(relic.ID()); got != card.House.Shadows {
				t.Errorf("artifact house = %v, want Shadows", got)
			}
		},
	)

	t.Run("keeps an on-identity artifact's house unchanged", func(t *testing.T) {
		var relic ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, Hand: ct.Cards(Sneklifter)},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&relic, ct.Artifact(ct.OfHouse(card.House.Mars)))),
			},
		})
		h.Game().
			SetPlayerHouses(0, []engine.House{card.House.Shadows, card.House.Mars, card.House.Sanctum})

		h.P1.Play(Sneklifter)

		if !slices.Contains(h.Game().Artifacts(0), relic.ID()) {
			t.Error("the artifact should be controlled by P1")
		}
		if got := h.Game().House(relic.ID()); got != card.House.Mars {
			t.Errorf("artifact house = %v, want Mars", got)
		}
	})
}
