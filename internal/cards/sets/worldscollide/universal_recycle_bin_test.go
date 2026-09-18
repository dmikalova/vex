package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Universal Recycle Bin
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Archive a card from your purge pile.
func TestUniversalRecycleBin(t *testing.T) {
	t.Run("archives a purged card you own", func(t *testing.T) {
		var recycled ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:   card.House.Logos,
				InPlay:  ct.Cards(UniversalRecycleBin),
				Discard: ct.Cards(ct.Bind(&recycled, ct.Creature())),
			},
		})

		// Set the card aside in the purge pile, then recycle it back to archives.
		h.Game().PurgeFromDiscard(0, recycled.ID())
		h.Expect(recycled).At(ct.Purge)

		h.P1.UseAction(UniversalRecycleBin)

		h.Expect(recycled).At(ct.Archives)
	})
}
