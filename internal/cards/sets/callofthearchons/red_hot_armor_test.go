package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

func TestRedHotArmor(t *testing.T) {
	var armored, bare ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Dis,
			Hand:  ct.Cards(RedHotArmor),
		},
		P2: ct.Side{
			InPlay: ct.Cards(
				ct.Bind(&armored, StaunchKnight),
				ct.Bind(&bare, Bumpsy),
			),
		},
	})

	h.P1.Play(RedHotArmor)

	// Staunch Knight has 2 armor, so it loses both points and takes 2 damage.
	h.Expect(armored).Damage(2)
	h.Expect(bare).Damage(0)
}
