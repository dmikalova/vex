package worldscollide

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Nizak, The Forgotten
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Dragon • Psion
//
//	While fighting, Nizak, The Forgotten gains invulnerable.
//	After a creature is destroyed in a fight with Nizak, The Forgotten, put it into its owner's hand.
func TestNizakTheForgotten(t *testing.T) {
	t.Run("an enemy that fights it dies and returns to its owner's hand", func(t *testing.T) {
		var nizak, attacker ct.Card
		h := ct.Play(t, ct.Setup{
			// The active player attacks Nizak, which the other player controls, so the
			// reaction fires on Nizak (the enemy of the destroyed attacker).
			P1: ct.Side{
				House: card.House.Brobnar,
				InPlay: ct.Cards(
					ct.Bind(&attacker, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(6))),
				),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&nizak, NizakTheForgotten))},
		})
		attacker.Ready()

		h.P1.Fight(attacker, nizak)

		// Invulnerable while fighting: equal power would trade both, but Nizak takes no
		// damage and survives.
		h.Expect(nizak).At(ct.PlayArea).Damage(0)
		// The attacker died to fightback and returned to its owner's (P1's) hand.
		h.Expect(attacker).At(ct.Hand)
	})

	t.Run("it destroys an enemy it fights and returns it to its owner's hand", func(t *testing.T) {
		var nizak, prey ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&nizak, NizakTheForgotten)),
			},
			P2: ct.Side{InPlay: ct.Cards(ct.Bind(&prey, ct.Creature(ct.Power(3))))},
		})
		nizak.Ready()

		h.P1.Fight(nizak, prey)

		// Nizak destroys the 3-power prey and, while fighting, takes no fightback.
		h.Expect(nizak).At(ct.PlayArea).Damage(0)
		// The destroyed enemy returned to its owner's (P2's) hand, not the discard.
		h.Expect(prey).At(ct.Hand)
	})

	t.Run("outside combat it is not invulnerable", func(t *testing.T) {
		var nizak ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Brobnar,
				InPlay: ct.Cards(ct.Bind(&nizak, NizakTheForgotten)),
			},
		})

		// A direct destruction outside a fight is not refused: the grant is
		// combat-scoped, so Nizak is destroyed and goes to the discard.
		h.Game().DestroyEach(0, []engine.LocalID{nizak.ID()})

		h.Expect(nizak).At(ct.Discard)
	})
}
