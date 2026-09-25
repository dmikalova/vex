package deckgen

import "github.com/dmikalova/vex/internal/engine"

// enhanceCap is the most bonus icons a card may carry, counting its printed icons
// and any Enhance icons landed on it (docs/deck-generation.md).
const enhanceCap = 5

// applyEnhancements is the deck-wide Enhance finishing pass (ADR 0004): after the
// deck is valid, every Enhance source's icons are distributed onto random cards in
// the same deck. A source contributes its icons and keeps none; each icon lands on
// a uniformly random eligible card — any card, the source included — capped at
// enhanceCap icons per card (printed icons count) and never of a kind a card
// barred with WithoutEnhancement. It draws from the generator's seeded RNG, so the
// distribution is deterministic for a given (Set, seed).
func (g *generator) applyEnhancements(deck *Deck) {
	slots := make([]*Slot, 0, DeckSize)
	for p := range deck.Pods {
		for s := range deck.Pods[p].Slots {
			slots = append(slots, &deck.Pods[p].Slots[s])
		}
	}
	// Clone each card's printed icons before landing anything: materialize copies
	// defs shallowly, so two slots of the same card share one Bonuses backing array
	// until this pass gives each its own.
	for _, s := range slots {
		bonusCopy := make([]engine.BonusIcon, len(s.Card.Bonuses))
		copy(bonusCopy, s.Card.Bonuses)
		s.Card.Bonuses = bonusCopy
	}
	var icons []engine.BonusIcon
	for _, s := range slots {
		icons = append(icons, s.Card.Enhances...)
	}
	if len(icons) == 0 {
		return
	}
	// Eligibility is per icon: a card is a candidate for an icon when it is under
	// the cap and does not bar that icon's kind, so a card that bars one kind can
	// still receive others.
	eligible := make([]*Slot, 0, len(slots))
	for _, ic := range icons {
		eligible = eligible[:0]
		for _, s := range slots {
			if len(s.Card.Bonuses) < enhanceCap && !s.Card.BarsEnhanceIcon(ic) {
				eligible = append(eligible, s)
			}
		}
		if len(eligible) == 0 {
			continue
		}
		s := eligible[g.r.Intn(len(eligible))]
		s.Card.Bonuses = append(s.Card.Bonuses, ic)
	}
}
