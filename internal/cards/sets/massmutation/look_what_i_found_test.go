package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Look What I Found!
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//
//	Omega.
//	Play: Put a tactic, artifact, creature, and upgrade from your discard pile into your hand.
func TestLookWhatIFound(t *testing.T) {
	t.Run("returns one card of each type from discard to hand", func(t *testing.T) {
		var tactic, artifact, creature, upgrade ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Untamed,
				Hand:  ct.Cards(LookWhatIFound),
				Discard: ct.Cards(
					ct.Bind(&tactic, ct.Tactic()),
					ct.Bind(&artifact, ct.Artifact()),
					ct.Bind(&creature, ct.Creature()),
					ct.Bind(&upgrade, ct.Upgrade()),
				),
			},
		})

		h.P1.Play(LookWhatIFound)

		h.Expect(tactic).At(ct.Hand)
		h.Expect(artifact).At(ct.Hand)
		h.Expect(creature).At(ct.Hand)
		h.Expect(upgrade).At(ct.Hand)
	})
}
