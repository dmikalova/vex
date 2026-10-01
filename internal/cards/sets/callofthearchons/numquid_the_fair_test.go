package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vex/internal/card"
	ct "github.com/dmikalova/vex/internal/cards/cardtest"
)

// Numquid the Fair
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human
//
//	Play: Destroy an enemy creature. If you are overwhelmed, repeat this effect.
func TestNumquidTheFair(t *testing.T) {
	t.Run("destroys enemy creatures while overwhelmed", func(t *testing.T) {
		var foe1, foe2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(NumquidTheFair),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&foe1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				ct.Bind(&foe2, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})

		// Numquid enters: P1 controls 1, P2 controls 2 -> overwhelmed. One destroy
		// brings it to parity, so the loop stops after a single destruction.
		h.P1.Play(NumquidTheFair)
		h.P1.ClickCard(foe1)

		h.Expect(foe1).At(ct.Discard)
		h.Expect(foe2).At(ct.PlayArea)
	})

	t.Run("a warded creature does not stop the repeat", func(t *testing.T) {
		// The repeat turns on a fact about the board, not on the destroy landing:
		// the ward absorbs the first destruction and the loop comes back for the
		// next creature while P1 is still overwhelmed.
		var warded, foe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(NumquidTheFair),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&warded, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				ct.Bind(&foe, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
			)},
		})
		warded.Ward()

		// Numquid enters: P1 controls 1, P2 controls 2 -> overwhelmed. The ward
		// absorbs the first destroy and the creature stays, so P1 is still
		// overwhelmed and the loop comes back for the second creature.
		h.P1.Play(NumquidTheFair)
		h.P1.ClickCard(warded)
		h.P1.ClickCard(foe)

		h.Expect(warded).At(ct.PlayArea)
		if warded.Warded() {
			t.Error("the ward should have been spent absorbing the destroy")
		}
		h.Expect(foe).At(ct.Discard)
	})

	// The same ruling the Neutron Shark test quotes: once the source is no longer
	// in play, its ability cannot be resolved again, even while the gate's
	// condition still holds (RAW §190; the pass under way finishes, RAW §191).
	t.Run("stops once Numquid itself has left play", func(t *testing.T) {
		var host, foe1, foe2 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				Hand:  ct.Cards(NumquidTheFair),
			},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Upgraded(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
					PhoenixHeart,
				),
				ct.Bind(&foe1, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
				ct.Bind(&foe2, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(5))),
			)},
		})

		// Destroying the Phoenix Heart host deals 3 damage to each creature, which
		// kills Numquid (power 3) mid-loop. P1 controls 0 creatures to P2's 2, so
		// the overwhelmed gate still holds — but the source is gone, so the two
		// surviving enemy creatures are safe.
		h.P1.Play(NumquidTheFair)
		h.P1.ClickCard(host)

		h.Expect(host).At(ct.Hand)
		h.Expect(NumquidTheFair).At(ct.Discard)
		h.Expect(foe1).At(ct.PlayArea).Damage(3)
		h.Expect(foe2).At(ct.PlayArea).Damage(3)
	})
}
