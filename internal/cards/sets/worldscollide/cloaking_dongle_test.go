package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Cloaking Dongle
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature and each of its neighbors gains elusive.
func TestCloakingDongle(t *testing.T) {
	t.Run("grants elusive to its host and both neighbors", func(t *testing.T) {
		var far, left, host, right ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Bind(&far, ct.Creature()),
					ct.Bind(&left, ct.Creature()),
					ct.Upgraded(ct.Bind(&host, ct.Creature()), CloakingDongle),
					ct.Bind(&right, ct.Creature()),
				),
			},
		})

		if !h.Game().HasKeyword(host.ID(), card.Keyword.Elusive) {
			t.Error("the host should gain elusive")
		}
		if !h.Game().HasKeyword(left.ID(), card.Keyword.Elusive) {
			t.Error("the left neighbor should gain elusive")
		}
		if !h.Game().HasKeyword(right.ID(), card.Keyword.Elusive) {
			t.Error("the right neighbor should gain elusive")
		}
		if h.Game().HasKeyword(far.ID(), card.Keyword.Elusive) {
			t.Error("a non-neighbor should not gain elusive")
		}
	})
}
