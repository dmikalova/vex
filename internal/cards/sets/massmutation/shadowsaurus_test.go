package massmutation

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Shadowsaurus
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Thief
//
//	Action: Move all Æmber from an enemy creature to your opponent's pool. If there was any Æmber on that creature, take control of it, and it belongs to house Shadows.
func TestShadowsaurus(t *testing.T) {
	t.Run("moves the Æmber and seizes the emptied creature as Shadows",
		func(t *testing.T) {
			var saurus, foe ct.Card
			h := ct.Play(t, ct.Setup{
				P1: ct.Side{
					House:  card.House.Shadows,
					InPlay: ct.Cards(ct.Bind(&saurus, Shadowsaurus)),
				},
				P2: ct.Side{InPlay: ct.Cards(
					ct.Bind(&foe, ct.Creature(ct.Power(3))),
				)},
			})
			h.Game().AddAmberOn(foe.ID(), 2)

			h.P1.UseAction(Shadowsaurus)
			h.P1.ClickOption("right flank")

			h.P2.ExpectAmber(2)
			if got := foe.AmberOn(); got != 0 {
				t.Errorf("foe Æmber = %d, want 0", got)
			}
			if !slices.Contains(h.Game().Battleline(0), foe.ID()) {
				t.Error("the emptied creature should be controlled by P1")
			}
			if got := h.Game().House(foe.ID()); got != card.House.Shadows {
				t.Errorf("seized creature house = %v, want Shadows", got)
			}
		})

	t.Run("without Æmber it does not take control", func(t *testing.T) {
		var saurus, bare ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Shadows,
				InPlay: ct.Cards(ct.Bind(&saurus, Shadowsaurus)),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&bare, ct.Creature(ct.Power(3))),
			)},
		})

		h.P1.UseAction(Shadowsaurus)

		if slices.Contains(h.Game().Battleline(0), bare.ID()) {
			t.Error("the empty creature should stay under P2's control")
		}
		if got := h.Game().House(bare.ID()); got == card.House.Shadows {
			t.Error("an untaken creature should keep its own house")
		}
	})
}
