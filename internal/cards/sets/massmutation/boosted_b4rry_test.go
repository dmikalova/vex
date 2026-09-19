package massmutation

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Boosted B4-RRY
//
//	House:  Shadows
//	Type:   Gigantic Creature
//	Rarity: Special
//	Power:  7
//	Armor:  2
//	Traits: Robot
//
//	Play/Fight/Reap: Choose one:
//	- Take control of an enemy artifact. If it does not belong to a house on your identity, it belongs to house Shadows.
//	- Play a random card from your opponent's archives.
func TestBoostedB4RRY(t *testing.T) {
	t.Run("takes control of an off-identity artifact and makes it Shadows",
		func(t *testing.T) {
			var b4rry, relic ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Shadows,
					InPlay: ct.Cards(ct.Bind(&b4rry, BoostedB4RRY)),
				},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&relic, ct.Artifact(ct.OfHouse(card.House.Untamed))),
				)},
			})
			h.Game().SetPlayerHouses(0, []engine.House{
				card.House.Shadows, card.House.Logos, card.House.Sanctum,
			})

			h.P1.Reap(b4rry)
			h.P1.ClickOption("take control")

			if !slices.Contains(h.Game().Artifacts(0), relic.ID()) {
				t.Error("the artifact should be controlled by P1")
			}
			if got := h.Game().House(relic.ID()); got != card.House.Shadows {
				t.Errorf("artifact house = %v, want Shadows", got)
			}
		})

	t.Run("plays a random card from the opponent's archives as yours",
		func(t *testing.T) {
			var b4rry, foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Shadows,
					InPlay: ct.Cards(ct.Bind(&b4rry, BoostedB4RRY)),
				},
				P2: ct.Side{Archives: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(3))),
				)},
			})

			h.P1.Reap(b4rry)
			h.P1.ClickOption("random card")

			if !h.Game().InPlay(foe.ID()) {
				t.Fatal("the archived creature should be in play")
			}
			if got := h.Game().Controller(foe.ID()); got != 0 {
				t.Errorf("controller = %d, want 0 (played as yours)", got)
			}
			if got := h.Game().Owner(foe.ID()); got != 1 {
				t.Errorf("owner = %d, want 1 (unchanged)", got)
			}
		})
}
