package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Master of 1
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with power 1.
func TestMasterOf1(t *testing.T) { testMaster(t, MasterOf1, 1) }

// Master of 2
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with power 2.
func TestMasterOf2(t *testing.T) { testMaster(t, MasterOf2, 2) }

// Master of 3
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with power 3.
func TestMasterOf3(t *testing.T) { testMaster(t, MasterOf3, 3) }

// Master of 4
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with power 4.
func TestMasterOf4(t *testing.T) { testMaster(t, MasterOf4, 4) }

// Master of 5
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with power 5.
func TestMasterOf5(t *testing.T) { testMaster(t, MasterOf5, 5) }

// testMaster exercises one Master of N variant: its Reap destroys only a
// creature of power n, and the destroy is optional — the controller may decline
// and still gain Æmber for reaping.
func testMaster(t *testing.T, variant card.Definition, n int) {
	t.Helper()

	t.Run("destroys only a creature of its own power", func(t *testing.T) {
		var master, matching, mismatched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&master, variant)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&matching, ct.Creature(ct.Power(n))),
					ct.Bind(&mismatched, ct.Creature(ct.Power(n+1))),
				),
			},
		})

		h.P1.Reap(master)
		h.P1.ClickCard(matching)

		h.Expect(matching).At(ct.Discard)
		h.Expect(mismatched).At(ct.PlayArea)
	})

	t.Run("may decline to destroy and still gains Æmber", func(t *testing.T) {
		var master, spared ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&master, variant)),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&spared, ct.Creature(ct.Power(n)))),
			},
		})

		h.P1.Reap(master)
		h.P1.ClickDone()

		h.Expect(spared).At(ct.PlayArea)
		h.P1.ExpectAmber(1)
	})
}
