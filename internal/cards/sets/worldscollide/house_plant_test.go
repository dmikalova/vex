package worldscollide

import (
	"math/rand"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// House Plant
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Elf • Thief
//
//	Elusive.
//
//	Template: its concrete card is materialized per deck at generation.
func TestHousePlant(t *testing.T) {
	deckWith := func(partner engine.House) card.SlotContext {
		return card.SlotContext{
			House:      card.House.Shadows,
			DeckHouses: [3]engine.House{engine.Shadows, partner, engine.HouseNone},
		}
	}

	t.Run("gains 1 Æmber when a player chooses the partner House", func(t *testing.T) {
		def := plantFor(deckWith(engine.Brobnar), rand.New(rand.NewSource(1)))
		if def.Name != "Brobnar Plant" {
			t.Fatalf("name = %q, want Brobnar Plant", def.Name)
		}

		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(def)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Brobnar)

		h.P1.ExpectAmber(1)
	})

	t.Run("does nothing when a different house is chosen", func(t *testing.T) {
		def := plantFor(deckWith(engine.Brobnar), rand.New(rand.NewSource(1)))

		h := ct.Play(t, ct.Setup{
			P1: ct.Side{House: card.House.Shadows, InPlay: ct.Cards(def)},
		})

		h.P1.EndTurn()
		h.P2.ChooseHouse(card.House.Logos)

		h.P1.ExpectAmber(0)
	})

	t.Run("materializes a Plant for a House WotC never printed one for", func(t *testing.T) {
		// WotC printed Plants for six Houses but not Sanctum; the template covers
		// it with no new card — the point of keying on the House enum.
		def := plantFor(deckWith(engine.Sanctum), rand.New(rand.NewSource(1)))
		if def.Name != "Sanctum Plant" {
			t.Fatalf("name = %q, want Sanctum Plant", def.Name)
		}
	})

	t.Run("falls back to a random partner with no deck Houses", func(t *testing.T) {
		// The uniqueness sampling test materializes with no DeckHouses set; the
		// fallback still binds to some non-Shadows House so every variant is reachable.
		def := plantFor(card.SlotContext{}, rand.New(rand.NewSource(1)))
		if def.House != card.House.Shadows {
			t.Fatalf("house = %v, want Shadows", def.House)
		}
		if def.Name == "Shadows Plant" {
			t.Fatal("a Plant never watches for its own Shadows")
		}
	})
}
