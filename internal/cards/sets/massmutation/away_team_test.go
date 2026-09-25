package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
	"github.com/dmikalova/vex/internal/engine"
)

// Away Team
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Alien • Human • Robot
//
//	Destroyed: Archive each upgrade on Away Team from play.
func TestAwayTeam(t *testing.T) {
	t.Run("archives its upgrades when it is destroyed", func(t *testing.T) {
		var awayTeam, up1, up2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.StarAlliance,
				InPlay: ct.Cards(
					ct.Upgraded(
						ct.Bind(&awayTeam, AwayTeam),
						ct.Bind(&up1, ct.Upgrade(ct.PowerBonus(1))),
						ct.Bind(&up2, ct.Upgrade(ct.PowerBonus(1))),
					),
				),
			},
		})

		h.Game().DestroyEach(0, []engine.LocalID{awayTeam.ID()})

		h.Expect(up1).At(ct.Archives)
		h.Expect(up2).At(ct.Archives)
	})
}
