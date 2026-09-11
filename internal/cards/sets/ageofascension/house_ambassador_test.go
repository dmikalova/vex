package ageofascension

import (
	"math/rand"
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// House Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//
//	Template: its concrete card is materialized per deck at generation.
func TestHouseAmbassador(t *testing.T) {
	deckWith := func(partner engine.House) card.SlotContext {
		return card.SlotContext{
			House:      card.House.Sanctum,
			DeckHouses: [3]engine.House{engine.Sanctum, partner, engine.HouseNone},
		}
	}

	t.Run("reap lets you play and use the partner House's cards", func(t *testing.T) {
		def := ambassadorFor(deckWith(engine.Mars), rand.New(rand.NewSource(1)))
		if def.Name != "Mars Ambassador" {
			t.Fatalf("name = %q, want Mars Ambassador", def.Name)
		}

		var marsInPlay ct.Card
		h := ct.Play(t, ct.Setup{
			P1: ct.Side{
				House: card.House.Sanctum,
				InPlay: ct.Cards(
					def,
					ct.Bind(&marsInPlay, ct.Creature(ct.OfHouse(card.House.Mars), ct.Power(3))),
				),
			},
			P2: ct.Side{},
		})

		h.P1.Reap(def)

		if got := h.Game().State.MayPlayHouse[0]; got != engine.Mars {
			t.Fatalf("MayPlayHouse = %v, want Mars", got)
		}
	})

	t.Run("materializes an Ambassador for a House it never printed", func(t *testing.T) {
		// Saurian and Star Alliance came after Age of Ascension, yet the template
		// covers them with no new card — the point of keying on the House enum.
		for _, partner := range []engine.House{engine.Saurian, engine.StarAlliance} {
			def := ambassadorFor(deckWith(partner), rand.New(rand.NewSource(1)))
			if want := partner.String() + " Ambassador"; def.Name != want {
				t.Fatalf("name = %q, want %q", def.Name, want)
			}
		}
	})

	t.Run("falls back to a random partner with no deck Houses", func(t *testing.T) {
		// The uniqueness sampling test materializes with no DeckHouses set; the
		// fallback still binds to some non-Sanctum House so every variant is reachable.
		def := ambassadorFor(card.SlotContext{}, rand.New(rand.NewSource(1)))
		if def.House != card.House.Sanctum {
			t.Fatalf("house = %v, want Sanctum", def.House)
		}
		if def.Name == "Sanctum Ambassador" {
			t.Fatal("an Ambassador never vouches for its own Sanctum")
		}
	})
}
