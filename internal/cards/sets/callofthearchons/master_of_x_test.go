package callofthearchons

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/deckgen"
)

// Master of X
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with power X.
func TestMasterOfX(t *testing.T) {
	t.Run("a variant destroys only a creature of its own power", func(t *testing.T) {
		var master, matching, mismatched ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&master, masterOfVariant(2))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(
					ct.Bind(&matching, ct.Creature(ct.Power(2))),
					ct.Bind(&mismatched, ct.Creature(ct.Power(3))),
				),
			},
		})

		h.P1.Reap(master)
		h.P1.ClickCard(matching)

		h.Expect(matching).At(ct.Discard)
		h.Expect(mismatched).At(ct.PlayArea)
	})

	t.Run("a variant with no creature of its power may destroy nothing", func(t *testing.T) {
		var master, safe ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House:  card.House.Dis,
				InPlay: ct.Cards(ct.Bind(&master, masterOfVariant(5))),
			},
			P2: ct.Side{
				InPlay: ct.Cards(ct.Bind(&safe, ct.Creature(ct.Power(3)))),
			},
		})

		h.P1.Reap(master)

		h.Expect(safe).At(ct.PlayArea)
		h.P1.ExpectAmber(1)
	})

	t.Run("materializes as one of the five numbered variants", func(t *testing.T) {
		m := masterOfXMaterializer(t)
		r := rand.New(rand.NewSource(1))
		seen := map[string]bool{}
		for range 200 {
			def := m.Materialize(deckgen.SlotContext{
				House:  card.House.Dis,
				Rarity: card.Rarity.Rare,
			}, r)
			seen[def.Name] = true
		}

		for n := 1; n <= 5; n++ {
			if want := fmt.Sprintf("Master of %d", n); !seen[want] {
				t.Errorf("%q never materialized", want)
			}
		}
		if len(seen) != 5 {
			t.Errorf(
				"materialized %d distinct names, want exactly the 5 variants: %v",
				len(seen),
				seen,
			)
		}
	})
}

// masterOfVariant builds one concrete Master of N the way the template's
// materializer does, without going through deck generation.
func masterOfVariant(n int) card.Definition {
	return card.Build(
		fmt.Sprintf("Master of %d", n),
		card.House.Dis,
		card.Type.Creature,
		card.Rarity.Rare,
		masterOf(n)...,
	)
}

// masterOfXMaterializer finds the template face's materializer in the registry.
func masterOfXMaterializer(t *testing.T) deckgen.Materializer {
	t.Helper()
	for _, rc := range card.Cards() {
		if rc.Def.Name == MasterOfX.Name {
			if rc.Materializer == nil {
				t.Fatal("Master of X registered without a materializer")
			}
			return rc.Materializer
		}
	}
	t.Fatal("Master of X is not registered")
	return nil
}
