package massmutation

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Sirs Colossus
//
//	House:  Sanctum
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  10
//	Armor:  3
//	Traits: Knight • Spirit
//
//	Taunt.
//	Play: Capture all your opponent's Æmber, distributed among any number of friendly creatures.
//	Fight: Move each Æmber on a friendly creature to the common supply.
func TestSirsColossus(t *testing.T) {
	var colossus ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House: card.House.Sanctum,
			Hand:  ct.Cards(ct.Bind(&colossus, SirsColossus), card.GiganticArt(SirsColossus)),
		},
		P2: ct.Side{Amber: 4},
	})

	// Play captures all the opponent's Æmber onto the sole friendly creature.
	h.P1.Play(colossus)
	h.P2.ExpectAmber(0)
	h.Expect(colossus).AmberOn(4)
}
