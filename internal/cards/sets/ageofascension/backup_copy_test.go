package ageofascension

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Backup Copy
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains, "Destroyed: Put this creature on top of its owner's deck."
func TestBackupCopy(t *testing.T) {
	t.Run("returns its host to the top of the deck when the host is destroyed", func(t *testing.T) {
		var host, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Logos,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Logos), ct.Power(2))),
						BackupCopy,
					),
				),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&foe, ct.Creature(ct.Power(6)))),
			},
		})

		h.P1.Fight(host, foe)

		h.Expect(host).At(ct.Deck)
	})
}
