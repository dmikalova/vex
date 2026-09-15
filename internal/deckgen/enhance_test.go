package deckgen

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/engine"
)

// fillDeck lays cards across the deck's 36 slots in order, repeating the last card
// to fill the remaining slots so no slot is left with a zero-value definition.
func fillDeck(cards ...engine.CardDefinition) *Deck {
	d := &Deck{}
	for i := 0; i < DeckSize; i++ {
		c := cards[len(cards)-1]
		if i < len(cards) {
			c = cards[i]
		}
		d.Pods[i/PodSize].Slots[i%PodSize].Card = c
	}
	return d
}

func totalBonuses(d *Deck) int {
	n := 0
	for p := range d.Pods {
		for s := range d.Pods[p].Slots {
			n += len(d.Pods[p].Slots[s].Card.Bonuses)
		}
	}
	return n
}

func plainCard(name string) engine.CardDefinition {
	return engine.NewCard(name, engine.Brobnar, engine.Creature, engine.Common, engine.WithPower(3))
}

func TestApplyEnhancementsLandsEveryIcon(t *testing.T) {
	src := engine.NewCard(
		"Source",
		engine.Brobnar,
		engine.Creature,
		engine.Common,
		engine.WithPower(
			3,
		),
		engine.WithEnhance(engine.BonusDamage, engine.BonusDamage, engine.BonusDraw),
	)
	deck := fillDeck(src, plainCard("A"), plainCard("B"))
	gen(NewSet("S", []Card{mkCard("B", engine.Brobnar, engine.Common)}, DefaultTuning())).
		applyEnhancements(deck)
	if got := totalBonuses(deck); got != 3 {
		t.Fatalf("landed icons = %d, want 3", got)
	}
}

func TestApplyEnhancementsIsDeterministic(t *testing.T) {
	src := engine.NewCard("Source", engine.Brobnar, engine.Creature, engine.Common,
		engine.WithPower(3), engine.WithEnhance(engine.BonusDamage, engine.BonusAember))
	a := fillDeck(src, plainCard("A"))
	b := fillDeck(src, plainCard("A"))
	set := NewSet("S", []Card{mkCard("B", engine.Brobnar, engine.Common)}, DefaultTuning())
	gen(set).applyEnhancements(a)
	gen(set).applyEnhancements(b)
	for i := 0; i < DeckSize; i++ {
		pa := a.Pods[i/PodSize].Slots[i%PodSize].Card.Bonuses
		pb := b.Pods[i/PodSize].Slots[i%PodSize].Card.Bonuses
		if len(pa) != len(pb) {
			t.Fatalf("slot %d differs: %v vs %v", i, pa, pb)
		}
	}
}

func TestApplyEnhancementsSkipsBarredKindAndRespectsCap(t *testing.T) {
	// One source floods far more Æmber icons than the deck can hold.
	flood := make([]engine.BonusIcon, 200)
	for i := range flood {
		flood[i] = engine.BonusAember
	}
	src := engine.NewCard("Source", engine.Brobnar, engine.Creature, engine.Common,
		engine.WithPower(3), engine.WithEnhance(flood...))
	barsAember := engine.NewCard("NoAember", engine.Logos, engine.Creature, engine.Common,
		engine.WithPower(3), engine.WithoutEnhancement(engine.BonusAember))
	barsCapture := engine.NewCard("NoCapture", engine.Logos, engine.Creature, engine.Common,
		engine.WithPower(3), engine.WithoutEnhancement(engine.BonusCapture))
	deck := fillDeck(src, barsAember, barsCapture, plainCard("A"))
	gen(NewSet("S", []Card{mkCard("B", engine.Brobnar, engine.Common)}, DefaultTuning())).
		applyEnhancements(deck)
	// The card that bars Æmber never receives an Æmber icon.
	if got := len(deck.Pods[0].Slots[1].Card.Bonuses); got != 0 {
		t.Fatalf("card barring Æmber got %d icons, want 0", got)
	}
	// A card that bars only Capture still receives the flooded Æmber icons.
	if got := len(deck.Pods[0].Slots[2].Card.Bonuses); got == 0 {
		t.Fatal("card barring only Capture should still receive Æmber icons")
	}
	// Every slot is filled at most to the cap; nothing exceeds it.
	for p := range deck.Pods {
		for s := range deck.Pods[p].Slots {
			if n := len(deck.Pods[p].Slots[s].Card.Bonuses); n > enhanceCap {
				t.Fatalf("slot exceeded cap: %d", n)
			}
		}
	}
}

func TestApplyEnhancementsNoSourcesIsNoop(t *testing.T) {
	deck := fillDeck(plainCard("A"), plainCard("B"))
	gen(NewSet("S", []Card{mkCard("B", engine.Brobnar, engine.Common)}, DefaultTuning())).
		applyEnhancements(deck)
	if got := totalBonuses(deck); got != 0 {
		t.Fatalf("bonuses = %d, want 0", got)
	}
}

func TestApplyEnhancementsPrintedIconsCountTowardCap(t *testing.T) {
	// Every card already carries the cap in printed icons, so no icon can land.
	full := engine.NewCard("Full", engine.Brobnar, engine.Creature, engine.Common,
		engine.WithPower(3), engine.WithBonus(
			engine.BonusAember, engine.BonusAember, engine.BonusAember,
			engine.BonusAember, engine.BonusAember),
		engine.WithEnhance(engine.BonusDraw))
	deck := fillDeck(full)
	gen(NewSet("S", []Card{mkCard("B", engine.Brobnar, engine.Common)}, DefaultTuning())).
		applyEnhancements(deck)
	if got := totalBonuses(deck); got != DeckSize*enhanceCap {
		t.Fatalf("bonuses = %d, want %d (nothing landed)", got, DeckSize*enhanceCap)
	}
}
