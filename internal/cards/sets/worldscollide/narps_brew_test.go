package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Narp's Brew
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains +2 armor and taunt.
func TestNarpsBrew(t *testing.T) {
	t.Run("host gains +2 armor and Taunt", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar))),
						NarpsBrew,
					),
				),
			},
		})

		h.Expect(host).Armor(2)
		if !h.Game().HasKeyword(host.ID(), card.Keyword.Taunt) {
			t.Error("the host should gain Taunt")
		}
	})
}
