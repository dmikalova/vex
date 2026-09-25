package ageofascension

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// [REDACTED]
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: [redacted]
//
//	After you choose Logos as your active house, place 1 Æmber from the common
//	supply on [REDACTED]. If there are 4 or more Æmber on it, destroy [REDACTED],
//	and forge a key at no cost.
func TestRedacted(t *testing.T) {
	t.Run("hoards 1 Æmber each time you choose Logos", func(t *testing.T) {
		var redacted ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&redacted, REDACTED)),
			},
			P2: ct.Side{},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)

		h.Expect(redacted).AmberOn(1)
	})

	t.Run("purges itself and forges a key at four Æmber", func(t *testing.T) {
		var redacted ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Logos,
				InPlay: ct.Cards(ct.Bind(&redacted, REDACTED)),
			},
			P2: ct.Side{},
		})

		h.Game().State.ForgeCanonicalKeys(0, 0)
		h.Game().AddAmberOn(redacted.ID(), 3)

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)
		h.P2.EndTurn()
		h.P1.ChooseHouse(card.House.Logos)

		h.Expect(redacted).At(ct.Purge)
		h.P1.ExpectKeys(1)
	})
}
