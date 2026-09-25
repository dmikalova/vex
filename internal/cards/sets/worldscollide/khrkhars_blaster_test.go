package worldscollide

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Khrkhar's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Khrkhar's Blaster to Lieutenant Khrkhar -> ward Lieutenant Khrkhar."
func TestKhrkharsBlaster(t *testing.T) {
	t.Run("attaches to Khrkhar and wards it", func(t *testing.T) {
		var carrier, khrkhar ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(
							&carrier,
							ct.Creature(ct.OfHouse(card.House.StarAlliance), ct.Power(4)),
						),
						KhrkharsBlaster,
					),
					ct.Bind(&khrkhar, LieutenantKhrkhar),
				),
			},
		})

		h.P1.Reap(carrier)
		h.P1.ExpectPrompt("Choose one")
		h.P1.ClickOption("attach")

		if !h.Game().Warded(khrkhar.ID()) {
			t.Error("Lieutenant Khrkhar should be warded")
		}
	})
}
