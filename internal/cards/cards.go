// Package cards is the card database aggregator. It blank-imports every released
// set package so their cards self-register (each card calls card.Register from
// its own init), then exposes the assembled database to the rest of the program.
//
// It lives outside package engine on purpose: set packages import the engine (via
// the card facade), and the engine must not import them back. Adding a card is
// just adding a self-registering file to a set; adding a set is one blank import
// below — neither requires listing individual cards here.
package cards

import (
	"sort"
	"strings"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/provenance"
	// Blank-imported so each set's cards self-register through its package init.
	_ "github.com/dmikalova/vactrol/internal/cards/sets/ageofascension"
	// Blank-imported so each set's cards self-register through its package init.
	_ "github.com/dmikalova/vactrol/internal/cards/sets/callofthearchons"
	"github.com/dmikalova/vactrol/internal/deckgen"
)

// All returns every registered card across every imported set.
func All() []card.Definition {
	return card.Registered()
}

// setName is the display name of the source set a card belongs to, taken from its
// first provenance tag. A card with no provenance (a wholly original card) is
// grouped under the empty name.
func setName(rc card.RegisteredCard) string {
	if len(rc.Provenance) == 0 {
		return ""
	}
	return rc.Provenance[0].Set.Name
}

// deckgenCard adapts a registered card to a deckgen pool entry.
func deckgenCard(rc card.RegisteredCard) deckgen.Card {
	return deckgen.Card{Def: rc.Def, Profile: rc.Profile, Materializer: rc.Materializer}
}

// setOrder returns the source-set names in release order, so grouping is
// deterministic and the newest set is last.
func setOrder() []string {
	var order []string
	for _, s := range provenance.Sets() {
		order = append(order, s.Name)
	}
	return order
}

// bySet groups the registered cards by source-set name.
func bySet() map[string][]card.RegisteredCard {
	groups := map[string][]card.RegisteredCard{}
	for _, rc := range card.Cards() {
		groups[setName(rc)] = append(groups[setName(rc)], rc)
	}
	return groups
}

// DeckSets assembles one deck-generation Set per source set. Each set's own pool
// is the cards implemented in its package plus its reprints — cards implemented in
// an earlier set but printed in this one too (declared by a set's 0set.go and
// resolved here by name), which are full pool members the same as a newly
// implemented card. Each set also carries a legacy pool of every other set's
// cards, from which a slot occasionally draws a same-House card (that draw may
// land on a card this set also reprints). The sets come back in release order.
func DeckSets() []deckgen.Set {
	groups := bySet()
	reprints := reprintsBySet()
	var sets []deckgen.Set
	for _, name := range setOrder() {
		own := ownPool(groups[name], reprints[name])
		if len(own) == 0 {
			continue
		}
		var legacy []deckgen.Card
		for other, otherRegs := range groups {
			if other == name {
				continue
			}
			for _, rc := range otherRegs {
				legacy = append(legacy, deckgenCard(rc))
			}
		}
		sets = append(
			sets,
			deckgen.NewSet(name, own, deckgen.DefaultTuning()).WithLegacy(legacy),
		)
	}
	return sets
}

// ownPool is a set's deck-generation pool: its natively implemented cards followed
// by its reprints, skipping a reprint whose name a native card already fills so a
// card is pooled once.
func ownPool(native, reprints []card.RegisteredCard) []deckgen.Card {
	seen := make(map[string]bool, len(native))
	out := make([]deckgen.Card, 0, len(native)+len(reprints))
	for _, rc := range native {
		out = append(out, deckgenCard(rc))
		seen[rc.Def.Name] = true
	}
	for _, rc := range reprints {
		if seen[rc.Def.Name] {
			continue
		}
		out = append(out, deckgenCard(rc))
		seen[rc.Def.Name] = true
	}
	return out
}

// reprintsBySet resolves each set's reprint claims (its 0set.go entries) to the
// registered card that implements them, keyed by source-set name. A claim whose
// card is not implemented is skipped. The per-set slices are sorted by name so the
// pool is deterministic regardless of package init order.
func reprintsBySet() map[string][]card.RegisteredCard {
	byName := map[string]card.RegisteredCard{}
	for _, rc := range card.Cards() {
		byName[normalizeName(rc.Def.Name)] = rc
	}
	out := map[string][]card.RegisteredCard{}
	for _, rp := range card.ReprintRefs() {
		if rc, ok := byName[normalizeName(rp.Name)]; ok {
			out[rp.Set.Name] = append(out[rp.Set.Name], rc)
		}
	}
	for _, regs := range out {
		sort.Slice(regs, func(i, j int) bool { return regs[i].Def.Name < regs[j].Def.Name })
	}
	return out
}

// normalizeName folds a card name to a case- and space-insensitive key so a
// reprint claim matches the card that implements it regardless of catalog casing.
func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// DeckSet assembles the deck-generation Set to build decks from: the first
// released set, with the other sets wired in as its legacy pool. It is the
// single-set entry point callers that generate one deck reach for, and it stays
// on the fully-implemented base set so generated decks are always full.
func DeckSet() deckgen.Set {
	return DeckSets()[0]
}

// DeckSetNames returns the deck-generation set names in release order, for a
// front-end that lets each player choose which set to play.
func DeckSetNames() []string {
	sets := DeckSets()
	names := make([]string, 0, len(sets))
	for _, s := range sets {
		names = append(names, s.Name)
	}
	return names
}

// DeckSetNamed returns the deck-generation Set with the given name — each carries
// its cross-set legacy pool — reporting false when no set matches.
func DeckSetNamed(name string) (deckgen.Set, bool) {
	for _, s := range DeckSets() {
		if s.Name == name {
			return s, true
		}
	}
	return deckgen.Set{}, false
}

// Set is a named group of cards, used for reporting such as the Card statistics
// view.
type Set struct {
	Name  string
	Cards []card.Definition
}

// Sets returns the card database grouped by source set, in release order.
func Sets() []Set {
	groups := bySet()
	var out []Set
	for _, name := range setOrder() {
		regs := groups[name]
		if len(regs) == 0 {
			continue
		}
		defs := make([]card.Definition, 0, len(regs))
		for _, rc := range regs {
			defs = append(defs, rc.Def)
		}
		sort.Slice(defs, func(i, j int) bool {
			if defs[i].House != defs[j].House {
				return defs[i].House < defs[j].House
			}
			return defs[i].Name < defs[j].Name
		})
		out = append(out, Set{Name: name, Cards: defs})
	}
	return out
}
