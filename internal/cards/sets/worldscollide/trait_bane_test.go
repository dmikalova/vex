package worldscollide

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Trait Bane
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Template: its concrete card is materialized per deck at generation.
func TestTraitBane(t *testing.T) {
	t.Run("destroys a creature of each of the three Houses' top traits", func(t *testing.T) {
		h1, h2, h3 := engine.Brobnar, engine.Logos, engine.Shadows
		bane := baneForHouses(h1, h2, h3)

		var c1, c2, c3 ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Dis, Hand: ct.Cards(bane)},
			P2: ct.Side{InPlay: ct.Cards(
				ct.Bind(&c1, ct.Creature(ct.Traits(card.MostCommonCreatureTrait(h1)))),
				ct.Bind(&c2, ct.Creature(ct.Traits(card.MostCommonCreatureTrait(h2)))),
				ct.Bind(&c3, ct.Creature(ct.Traits(card.MostCommonCreatureTrait(h3)))),
			)},
		})

		h.P1.Play(bane)

		h.Expect(c1).At(ct.Discard)
		h.Expect(c2).At(ct.Discard)
		h.Expect(c3).At(ct.Discard)
		h.P1.ExpectAmber(1)
	})

	t.Run("splices a growing window of each trait into the name", func(t *testing.T) {
		got := baneName([]engine.Trait{card.Traits.Beast, card.Traits.Scientist, card.Traits.Thief})
		if got != "Bescithief's Bane" {
			t.Fatalf("baneName = %q, want Bescithief's Bane", got)
		}
	})

	t.Run("names are order-independent for the same three Houses", func(t *testing.T) {
		a := baneForHouses(engine.Brobnar, engine.Logos, engine.Shadows)
		b := baneForHouses(engine.Shadows, engine.Brobnar, engine.Logos)
		if a.Name != b.Name {
			t.Fatalf("names differ by House order: %q vs %q", a.Name, b.Name)
		}
	})

	t.Run("materializes a Dis bane from three random Houses", func(t *testing.T) {
		def := baneFor(card.SlotContext{}, rand.New(rand.NewSource(1)))
		if def.House != card.House.Dis {
			t.Fatalf("house = %v, want Dis", def.House)
		}
		if !strings.HasSuffix(def.Name, " Bane") {
			t.Fatalf("name = %q, want a ... Bane", def.Name)
		}
		if got := engine.RenderCardText(&def); !strings.Contains(got, "Destroy a ") {
			t.Fatalf("text = %q, want a destroy clause", got)
		}
	})
}

func TestWindow(t *testing.T) {
	for _, tc := range []struct {
		in   string
		n    int
		want string
	}{
		{"Beast", 2, "be"},
		{"Scientist", 3, "sci"},
		{"Thief", 4, "thief"},
		{"Demon", 2, "dem"},
		{"AI", 4, "ai"},
	} {
		if got := window(tc.in, tc.n); got != tc.want {
			t.Errorf("window(%q, %d) = %q, want %q", tc.in, tc.n, got, tc.want)
		}
	}
}
