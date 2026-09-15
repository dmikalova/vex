package callofthearchons

import (
	"slices"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Collar of Subordination
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Rare
//
//	Play: Take control of this creature until Collar of Subordination leaves play.
func TestCollarOfSubordination(t *testing.T) {
	t.Run("takes control of an enemy creature while the collar stays attached", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(CollarOfSubordination),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Dis), ct.Power(3))),
				),
			},
		})

		h.P1.Play(CollarOfSubordination)
		host.Ready()

		expectControlledBy(t, h.Game(), host.ID(), 0)
		h.P1.Reap(host)
		h.P1.ExpectAmber(1)
		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Dis)
		if err := h.Game().CanUse(1, host.ID()); err != engine.ErrWrongType {
			t.Fatalf("P2 CanUse(controlled host) = %v, want wrong type", err)
		}
	})

	t.Run("puts a destroyed controlled creature into its owner's discard", func(t *testing.T) {
		var host ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Dis,
				Hand:  ct.Cards(CollarOfSubordination),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&host, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Play(CollarOfSubordination)
		h.Game().DealDamage(0, []engine.DamageTarget{{ID: host.ID(), Amount: 3}})

		h.Expect(host).At(ct.Discard)
		if !slices.Contains(h.Game().Discard(1), host.ID()) {
			t.Fatalf(
				"controlled host went to discard %v/%v, want P2 discard",
				h.Game().Discard(0),
				h.Game().Discard(1),
			)
		}
	})
}

func expectControlledBy(t *testing.T, g *engine.Game, id engine.LocalID, player int) {
	t.Helper()
	if !slices.Contains(g.Battleline(player), id) {
		t.Fatalf(
			"%s is in battlelines %v/%v, want P%d",
			g.Name(id),
			g.Battleline(0),
			g.Battleline(1),
			player+1,
		)
	}
	if slices.Contains(g.Battleline(1-player), id) {
		t.Fatalf("%s is also in P%d battleline", g.Name(id), 2-player)
	}
}
