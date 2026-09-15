package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Whispering Reliquary
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Put an artifact into its owner's hand.
func TestWhisperingReliquary(t *testing.T) {
	t.Run("puts a chosen artifact into its owner's hand", func(t *testing.T) {
		var foeArtifact ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Sanctum, InPlay: ct.Cards(WhisperingReliquary)},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&foeArtifact, ct.Artifact()))},
		})

		h.P1.UseAction(WhisperingReliquary)
		h.P1.ClickCard(foeArtifact)

		h.Expect(foeArtifact).At(ct.Hand)
	})
}
