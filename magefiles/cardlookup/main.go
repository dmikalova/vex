// Command cardlookup queries the embedded KeyForge source catalogs (the
// provenance data under internal/cards/provenance) and the implemented card
// database (internal/cards). It exists so card-authoring work never needs a
// throwaway grep/JSON script to find a card's collector number, printed text, or
// implementation status.
//
// Subcommands (run via `mage lookup`, `mage missing`, `mage coverage`):
//
//	lookup <query>       Print every source card whose name contains <query>
//	                     (case-insensitive), with set code, collector number,
//	                     house, type, rarity, and printed text — everything a
//	                     card.New(...) definition needs, including the number for
//	                     card.Provenance(card.CotA, <n>).
//
//	missing [setSlug]    List the source cards in a set (default: the Call of the
//	                     Archons catalog) that no implemented card yet tags with a
//	                     provenance Ref, i.e. the cards still to implement.
//
//	coverage             Print, per source set, how many of its cards are covered
//	                     by an implemented card's provenance Ref.
//
//	stub <setSlug>       Generate a build-excluded (`//go:build todo`) stub file
//	                     for every unimplemented card in a set, each carrying the
//	                     printed text and a TODO marker. Excluded stubs do not
//	                     compile or register, so the database and coverage stay
//	                     honest until a card is actually implemented.
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
		return coverage()
	case "stub":
		return stub(args[1:])
	default:
		return usage()
	}
}

func usage() error {
	return fmt.Errorf(
		"usage: cardlookup <lookup <query> | missing [setSlug] | coverage | stub <setSlug>>",
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
	fmt.Printf("%s #%d  %s\n", set.Code, c.Number, c.Name)
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
	fmt.Printf("  Provenance: card.Provenance(card.%s, %d)\n\n", provCodeVar(set), c.Number)
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
// an implemented card tags with a provenance Ref.
func coveredNumbers() map[string]map[int]bool {
	covered := map[string]map[int]bool{}
	for _, rc := range card.Cards() {
		for _, ref := range rc.Provenance {
			if covered[ref.Set.Slug] == nil {
				covered[ref.Set.Slug] = map[int]bool{}
			}
			covered[ref.Set.Slug][ref.Number] = true
		}
	}
	return covered
}

// missing lists the source cards in a set not yet covered by any implemented card.
func missing(args []string) error {
	slug := provenance.CallOfTheArchons.Slug
	if len(args) == 1 {
		slug = args[0]
	} else if len(args) > 1 {
		return fmt.Errorf("usage: cardlookup missing [setSlug]")
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

// coverage prints, per source set, how many cards are covered by an implemented
// card's provenance Ref.
func coverage() error {
	_ = cards.All()
	covered := coveredNumbers()
	for _, set := range provenance.Sets() {
		n := 0
		for _, c := range set.Cards {
			if covered[set.Slug][c.Number] {
				n++
			}
		}
		fmt.Printf("%-5s %-22s %3d / %3d\n", set.Code, set.Name, n, len(set.Cards))
	}
	return nil
}
