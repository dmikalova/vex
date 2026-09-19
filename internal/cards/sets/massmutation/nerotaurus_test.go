package massmutation

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
)

// Nerotaurus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Dinosaur • Politician
//
//	Fight: Your opponent cannot use creatures to reap during their next turn.
//	Reap: Your opponent cannot use creatures to fight during their next turn.
func TestNerotaurus(t *testing.T) {
	t.Run("fight bars enemy creatures from reaping next turn", func(t *testing.T) {
		var enemy ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(Nerotaurus),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&enemy, ct.Creature(ct.OfHouse(card.House.Brobnar), ct.Power(3))),
				),
			},
		})

		h.P1.Fight(Nerotaurus, enemy)

		if !h.Game().State.CannotReapNext[1].Value {
			t.Error("Nerotaurus fight should bar the opponent from reaping next turn")
		}
		if h.Game().State.CannotFightNext[0].Value {
			t.Error("Nerotaurus should not restrict the caster")
		}
	})

	t.Run("reap bars enemy creatures from fighting next turn", func(t *testing.T) {
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Saurian,
				InPlay: ct.Cards(Nerotaurus),
			},
		})

		h.P1.Reap(Nerotaurus)

		if !h.Game().State.CannotFightNext[1].Value {
			t.Error("Nerotaurus reap should bar the opponent from fighting next turn")
		}
		if h.Game().State.CannotFightNext[0].Value {
			t.Error("Nerotaurus should not restrict the caster")
		}
	})
}
