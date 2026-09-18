// Command cardlookup queries the embedded KeyForge source catalogs (the
// provenance data under internal/cards/provenance) and the implemented card
// database (internal/cards). It exists so card-authoring work never needs a
// throwaway grep/JSON script to find a card's collector number, printed text, or
// implementation status.
//
// Subcommands (run via `mage tool:lookup`, `mage tool:missing`, `mage tool:coverage`):
//
//	lookup <query>       Print every source card whose name contains <query>
//	                     (case-insensitive), with set code, collector number,
//	                     house, type, rarity, and printed text — everything a
//	                     card.New(...) definition needs, including the number for
//	                     card.Provenance(card.CotA, <n>).
//
//	missing [setSlug]    List the source cards in a set (default: the Call of the
//	                     Archons catalog) that no implemented card yet tags with a
//	                     provenance Ref, i.e. the cards still to implement. With no
//	                     set named, an interactive ↑/↓ picker chooses one.
//
//	coverage             Print, per source set, how many of its cards are covered
//	                     by an implemented card's provenance Ref.
//
//	stub <setSlug>       Generate a build-excluded (`//go:build todo`) stub file
//	                     for every unimplemented card in a set, each carrying the
//	                     printed text and a TODO marker. Excluded stubs do not
//	                     compile or register, so the database and coverage stay
//	                     honest until a card is actually implemented.
//
//	next-card [setSlug]  Print the next unimplemented card whose stub still carries
//	                     the `//go:build todo` constraint, in collector-number
//	                     order — the card to build next.
//
//	node-usage           Print every exported name on the card facade, grouped by
//	                     the category its declaration block documents, with how
//	                     many card definitions and sets use it.
//
//	review [-n=<count>]  Open a random batch of card files (default 10) in VS Code
//	                     for review, recording them so they are not picked again
//	                     until the whole pool has been seen, then cycling.
//
//	import-provenance [setSlug|all]
//	                     Rebuild a set's source catalog (…/provenance/<slug>.json)
//	                     from the Master Vault decks feed, ASCII-folding names and
//	                     text and expanding the amber/damage markup. With no set it
//	                     opens a picker offering every set plus "All sets".
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"

	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/cards/provenance"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// isInteractive reports whether stdin is a terminal, so the interactive set
// picker is only offered when a person can actually drive it. A char-device
// check alone is not enough — /dev/null is a char device too — so it asks the
// terminal library.
func isInteractive() bool {
	return term.IsTerminal(os.Stdin.Fd())
}

func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}
	switch args[0] {
	case "lookup":
		return lookup(args[1:])
	case "missing":
		return missing(args[1:])
	case "coverage":
		return coverage(args[1:])
	case "stub":
		return stub(args[1:])
	case "next-card":
		return nextCard(args[1:])
	case "node-usage":
		return nodeUsage(args[1:])
	case "review":
		return review(args[1:])
	case "import-provenance":
		return importProvenance(args[1:])
	default:
		return usage()
	}
}

func usage() error {
	return fmt.Errorf(
		"usage: cardlookup <lookup <query> | missing [setSlug] | " +
			"coverage [-new] | stub <setSlug> | next-card [setSlug] | " +
			"node-usage [-max=<n>] [-category=<substring>] | " +
			"review [-n=<count>] | " +
			"import-provenance [setSlug|all]>",
	)
}

// lookup prints every source card whose name contains the query substring.
func lookup(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: cardlookup lookup <query>")
	}
	query := strings.ToLower(strings.Join(args, " "))
	found := 0
	for _, set := range provenance.Sets() {
		for _, c := range set.Cards {
			if !strings.Contains(strings.ToLower(c.Name), query) {
				continue
			}
			found++
			printCard(set.SourceSet, c)
		}
	}
	if found == 0 {
		return fmt.Errorf("no source card matches %q", strings.Join(args, " "))
	}
	return nil
}

// printCard renders one source card as an aligned block, ending with a ready-made
// card.Provenance call for the definition.
func printCard(set provenance.SourceSet, c provenance.Card) {
	fmt.Printf("%s #%s  %s\n", set.Code, c.Number, c.Name)
	fmt.Printf("  House:  %s\n", c.House)
	fmt.Printf("  Type:   %s\n", c.Type)
	fmt.Printf("  Rarity: %s\n", c.Rarity)
	if strings.EqualFold(c.Type, "creature") {
		fmt.Printf("  Power:  %d\n", c.Power)
		if c.Armor > 0 {
			fmt.Printf("  Armor:  %d\n", c.Armor)
		}
	}
	if c.Amber > 0 {
		fmt.Printf("  Æmber:  %d\n", c.Amber)
	}
	if len(c.Traits) > 0 {
		fmt.Printf("  Traits: %s\n", strings.Join(c.Traits, " • "))
	}
	if len(c.Keywords) > 0 {
		fmt.Printf("  Keywords: %s\n", strings.Join(c.Keywords, " • "))
	}
	if c.Text != "" {
		fmt.Printf("  Text:   %s\n", strings.ReplaceAll(c.Text, "\n", "\n          "))
	}
	fmt.Printf("  Provenance: card.Provenance(card.%s, %q)\n\n", provCodeVar(set), c.Number)
}

// provCodeVar maps a source set to the card-facade variable that names it (card.CotA,
// card.AoA, ...), falling back to the set code when unknown.
func provCodeVar(set provenance.SourceSet) string {
	switch set.Slug {
	case provenance.CallOfTheArchons.Slug:
		return "CotA"
	default:
		return set.Code
	}
}

// coveredNumbers returns, per source-set slug, the set of collector numbers that
// count as implemented. A source printing is implemented when some implemented
// card carries a provenance Ref to any printing of that same card: the Ref names
// a source card, and every same-named printing across sets is the same card, so
// tagging one printing covers them all. This means a card need only tag the
// printing it was built from (e.g. Labwork's CotA #114 and #271) — its reprints
// in later sets are covered by name without a per-set Ref, and a card renamed
// away from its printed name still covers its printings through the Ref.
func coveredNumbers() map[string]map[string]bool {
	covered := map[string]map[string]bool{}
	mark := func(slug string, number string) {
		if covered[slug] == nil {
			covered[slug] = map[string]bool{}
		}
		covered[slug][number] = true
	}

	refName := sourceNameByRef()
	implemented := map[string]bool{}
	for _, rc := range card.Cards() {
		// A card counts as implemented under its own name (so a same-named printing
		// in any set's catalog is covered) and under every source name its provenance
		// Refs resolve to. The own-name entry matters when a Ref does not resolve — an
		// anomaly's A0x Ref, since the source catalog omits its anomalies.
		implemented[normalizeName(rc.Def.Name)] = true
		for _, ref := range rc.Provenance {
			if name, ok := refName[refKey(ref)]; ok {
				implemented[normalizeName(name)] = true
			}
		}
	}
	for _, set := range provenance.Sets() {
		for _, c := range set.Cards {
			if implemented[normalizeName(c.Name)] {
				mark(set.Slug, c.Number)
			}
		}
	}
	return covered
}

// sourceNameByRef maps every source printing (set + collector number) to its
// catalog card name, so a card's provenance Ref can be resolved to the name of
// the source card it was built from. Refs are keyed through refKey so a call-site
// number ("4") matches its zero-padded catalog form ("004").
func sourceNameByRef() map[provenance.Ref]string {
	out := map[provenance.Ref]string{}
	for _, set := range provenance.Sets() {
		for _, c := range set.Cards {
			out[refKey(provenance.Ref{Set: set.SourceSet, Number: c.Number})] = c.Name
		}
	}
	return out
}

// refKey canonicalizes a Ref's collector number so a call-site number and its
// catalog form resolve to the same key: an all-digit number drops its leading
// zeros ("004" and "4" both become "4"), a lettered reference number is left as
// is ("S01", "A21").
func refKey(ref provenance.Ref) provenance.Ref {
	ref.Number = normNumber(ref.Number)
	return ref
}

// normNumber canonicalizes a collector number: an all-digit number is reduced to
// its integer form ("004" -> "4"), and any number with a letter is returned
// unchanged ("S01", "A21", "P07").
func normNumber(s string) string {
	if n, err := strconv.Atoi(s); err == nil {
		return strconv.Itoa(n)
	}
	return s
}

// normalizeName folds a card name to a case- and space-insensitive key so an
// implemented card matches its reprints across sets.
func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// missing lists the source cards in a set not yet covered by any implemented card.
// With no set named it opens the interactive picker on a terminal, and falls back
// to Call of the Archons when stdin is not a terminal (a pipe or CI).
func missing(args []string) error {
	var slug string
	switch len(args) {
	case 0:
	case 1:
		slug = args[0]
	default:
		return fmt.Errorf("usage: cardlookup missing [setSlug]")
	}
	if slug == "" {
		if isInteractive() {
			picked, err := pickSet()
			if err != nil {
				return err
			}
			if picked == "" {
				return nil // cancelled
			}
			slug = picked
		} else {
			slug = provenance.CallOfTheArchons.Slug
		}
	}
	// Touch the aggregator so every set package registers its cards.
	_ = cards.All()
	covered := coveredNumbers()[slug]

	sets := provenance.Sets()
	var set *provenance.Set
	for i := range sets {
		if sets[i].Slug == slug {
			set = &sets[i]
			break
		}
	}
	if set == nil {
		return fmt.Errorf("unknown set slug %q", slug)
	}

	var miss []provenance.Card
	for _, c := range set.Cards {
		if !covered[c.Number] {
			miss = append(miss, c)
		}
	}
	sort.Slice(miss, func(i, j int) bool { return miss[i].Number < miss[j].Number })
	fmt.Printf("%s: %d of %d cards not yet implemented\n\n", set.Name, len(miss), len(set.Cards))
	for _, c := range miss {
		printCard(set.SourceSet, c)
	}
	return nil
}

// nextCard prints the next unimplemented card whose stub file still carries the
// `//go:build todo` constraint, walking the set's missing cards in collector-
// number order and stopping at the first one that has a stub on disk. It is the
// pick-the-next-card step of the implement-cards workflow: build the card it
// names, drop the build tag, and run it again for the next. With no set named it
// resolves the slug the way `missing` does (interactive picker, or Call of the
// Archons when stdin is not a terminal).
func nextCard(args []string) error {
	var slug string
	switch len(args) {
	case 0:
	case 1:
		slug = args[0]
	default:
		return fmt.Errorf("usage: cardlookup next-card [setSlug]")
	}
	if slug == "" {
		if isInteractive() {
			picked, err := pickSet()
			if err != nil {
				return err
			}
			if picked == "" {
				return nil // cancelled
			}
			slug = picked
		} else {
			slug = provenance.CallOfTheArchons.Slug
		}
	}
	// Touch the aggregator so every set package registers its cards.
	_ = cards.All()
	covered := coveredNumbers()[slug]

	sets := provenance.Sets()
	var set *provenance.Set
	for i := range sets {
		if sets[i].Slug == slug {
			set = &sets[i]
			break
		}
	}
	if set == nil {
		return fmt.Errorf("unknown set slug %q", slug)
	}

	var miss []provenance.Card
	for _, c := range set.Cards {
		if !covered[c.Number] {
			miss = append(miss, c)
		}
	}
	sort.Slice(miss, func(i, j int) bool { return miss[i].Number < miss[j].Number })

	dir := filepath.Join("internal", "cards", "sets", slug)
	for _, c := range miss {
		path := filepath.Join(dir, fileName(c.Name)+".go")
		if !hasBuildTodo(path) {
			continue
		}
		fmt.Printf("%s\n\n", path)
		printCard(set.SourceSet, c)
		return nil
	}
	return fmt.Errorf(
		"no //go:build todo stub found in %s — run `mage tool:stub %s` first",
		dir,
		slug,
	)
}

// hasBuildTodo reports whether the file at path exists and starts with the
// `//go:build todo` constraint that marks an unimplemented stub.
func hasBuildTodo(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(string(data)), "//go:build todo")
}

// coverage prints, per source set, how many cards are covered by an implemented
// card's provenance Ref. With -new it counts only the cards a set introduces,
// excluding the ones it reprints from an earlier set, so the total is the set's
// own new cards rather than its whole printing.
func coverage(args []string) error {
	newOnly := false
	for _, a := range args {
		switch a {
		case "-new", "--new":
			newOnly = true
		default:
			return fmt.Errorf("usage: cardlookup coverage [-new]")
		}
	}
	_ = cards.All()
	covered := coveredNumbers()
	var introduced map[string]map[string]bool
	if newOnly {
		introduced = newNumbers()
	}
	var sumCovered, sumTotal int
	for _, set := range provenance.Sets() {
		n, total := 0, 0
		for _, c := range set.Cards {
			if newOnly && !introduced[set.Slug][c.Number] {
				continue
			}
			total++
			if covered[set.Slug][c.Number] {
				n++
			}
		}
		sumCovered += n
		sumTotal += total
		fmt.Printf("%-5s %-22s %3d / %3d  (%3d left)\n", set.Code, set.Name, n, total, total-n)
	}
	fmt.Printf("%-5s %-22s %3d / %3d  (%3d left)\n",
		"", "TOTAL", sumCovered, sumTotal, sumTotal-sumCovered)
	return nil
}

// newNumbers returns, per source-set slug, the collector numbers of the cards a
// set introduces: a card whose name first appears in that set walking the sets in
// release order. A card reprinted from an earlier set is left out, so the count is
// the set's own new cards rather than everything it prints.
func newNumbers() map[string]map[string]bool {
	seen := map[string]bool{}
	introduced := map[string]map[string]bool{}
	for _, set := range provenance.Sets() {
		for _, c := range set.Cards {
			key := normalizeName(c.Name)
			if seen[key] {
				continue
			}
			seen[key] = true
			if introduced[set.Slug] == nil {
				introduced[set.Slug] = map[string]bool{}
			}
			introduced[set.Slug][c.Number] = true
		}
	}
	return introduced
}
