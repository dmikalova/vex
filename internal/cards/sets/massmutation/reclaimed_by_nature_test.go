package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Reclaimed by Nature
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Purge an artifact. Resolve that card's bonus icons.
func TestReclaimedByNature(t *testing.T) {
	t.Run("purges an artifact and resolves its bonus icons", func(t *testing.T) {
		var artifact ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(ReclaimedByNature),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&artifact, ct.Artifact(ct.Bonus(card.Bonus.Aember))),
				),
			},
		})

		before := h.P1.Amber()
		h.P1.Play(ReclaimedByNature)

		// Reclaimed's own Æmber bonus gains 1; the purged artifact's Æmber bonus
		// gains 1 more, resolved for the player who played Reclaimed.
		if got := h.P1.Amber() - before; got != 2 {
			t.Fatalf("aember gained = %d, want 2", got)
		}
		h.Expect(artifact).At(ct.Purge)
	})
}
