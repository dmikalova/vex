package ageofascension_test

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/card"
	ct "github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/cards/sets/ageofascension"
)

// Shard of Life
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item • Shard
//
//	Action: For each friendly Shard, shuffle a card from your discard pile into your deck.
func TestShardOfLife(t *testing.T) {
	var shard ct.Card
	h := ct.Play(t, ct.Setup{
		P1: ct.Side{
			House:   card.House.Untamed,
			InPlay:  ct.Cards(ct.Bind(&shard, ageofascension.ShardOfLife)),
			Discard: ct.Cards(ct.Creature(ct.Power(1))),
		},
	})

	if got := len(h.Game().Discard(0)); got != 1 {
		t.Fatalf("discard size before = %d, want 1", got)
	}

	// Shard of Life is itself a Shard, so its count is one: it shuffles the sole
	// discard card back into the deck.
	h.P1.UseAction(shard)

	if got := len(h.Game().Discard(0)); got != 0 {
		t.Errorf("discard size after = %d, want 0", got)
	}
}
